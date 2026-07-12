package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PartHandler struct {
	db *gorm.DB
}

func NewPartHandler(deps Deps) *PartHandler {
	return &PartHandler{db: deps.DB}
}

type partRequest struct {
	Type           int8             `json:"type"`
	CategoryID     int64            `json:"categoryId"`
	Title          string           `json:"title"`
	Brand          string           `json:"brand"`
	Model          string           `json:"model"`
	ConditionLevel int8             `json:"conditionLevel"`
	Price          *float64         `json:"price"`
	CityCode       string           `json:"cityCode"`
	CityName       string           `json:"cityName"`
	Description    string           `json:"description"`
	ContactPolicy  int8             `json:"contactPolicy"`
	Images         []string         `json:"images"`
	Fitments       []fitmentRequest `json:"fitments"`
}

type fitmentRequest struct {
	VehicleModelID int64  `json:"vehicleModelId"`
	Note           string `json:"note"`
}

type partDetail struct {
	model.Part
	Images   []model.PartImage   `json:"images"`
	Fitments []model.PartFitment `json:"fitments"`
}

func (h *PartHandler) List(c *gin.Context) {
	page, pageSize := pagination(c)
	query := h.db.Model(&model.Part{}).Where("parts.status = ?", statusVisible)

	if partType, ok := queryInt8(c, "type"); ok {
		query = query.Where("parts.type = ?", partType)
	}
	if categoryID, ok := queryInt64(c, "categoryId"); ok {
		query = query.Where("parts.category_id = ?", categoryID)
	}
	if cityCode := strings.TrimSpace(c.Query("cityCode")); cityCode != "" {
		query = query.Where("parts.city_code = ?", cityCode)
	}
	if conditionLevel, ok := queryInt8(c, "conditionLevel"); ok {
		query = query.Where("parts.condition_level = ?", conditionLevel)
	}
	if minPrice, ok := queryFloat(c, "minPrice"); ok {
		query = query.Where("parts.price >= ?", minPrice)
	}
	if maxPrice, ok := queryFloat(c, "maxPrice"); ok {
		query = query.Where("parts.price <= ?", maxPrice)
	}
	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("parts.title LIKE ? OR parts.description LIKE ? OR parts.brand LIKE ? OR parts.model LIKE ?", pattern, pattern, pattern, pattern)
	}
	if userID, ok := queryInt64(c, "userId"); ok {
		query = query.Where("parts.user_id = ?", userID)
	}
	if vehicleModelID, ok := queryInt64(c, "vehicleModelId"); ok {
		query = query.Joins("JOIN part_fitments ON part_fitments.part_id = parts.id").
			Where("part_fitments.vehicle_model_id = ?", vehicleModelID)
	}

	var total int64
	if err := query.Distinct("parts.id").Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var parts []model.Part
	if err := query.Select("parts.*").
		Group("parts.id").
		Order("parts.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&parts).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{
		"items":    parts,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (h *PartHandler) Get(c *gin.Context) {
	partID, okID := pathID(c, "partId")
	if !okID {
		return
	}

	var part model.Part
	if err := h.db.Where("id = ? AND status IN ?", partID, []int8{statusVisible, statusPendingReview}).First(&part).Error; err != nil {
		serverError(c, err)
		return
	}

	var images []model.PartImage
	if err := h.db.Where("part_id = ?", partID).Order("sort_order ASC, id ASC").Find(&images).Error; err != nil {
		serverError(c, err)
		return
	}

	var fitments []model.PartFitment
	if err := h.db.Where("part_id = ?", partID).Order("id ASC").Find(&fitments).Error; err != nil {
		serverError(c, err)
		return
	}

	_ = h.db.Model(&model.Part{}).Where("id = ?", partID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error

	ok(c, partDetail{Part: part, Images: images, Fitments: fitments})
}

func (h *PartHandler) Create(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var req partRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if err := validatePartRequest(req); err != nil {
		badRequest(c, err.Error())
		return
	}

	status, err := contentStatus(h.db, req.Title, req.Brand, req.Model, req.Description)
	if err != nil {
		serverError(c, err)
		return
	}

	part := model.Part{
		UserID:         userID,
		Type:           req.Type,
		CategoryID:     req.CategoryID,
		Title:          normalizeText(req.Title, 120),
		Brand:          normalizeText(req.Brand, 64),
		Model:          normalizeText(req.Model, 128),
		ConditionLevel: req.ConditionLevel,
		Price:          req.Price,
		CityCode:       normalizeText(req.CityCode, 32),
		CityName:       normalizeText(req.CityName, 64),
		Description:    strings.TrimSpace(req.Description),
		ContactPolicy:  contactPolicyOrDefault(req.ContactPolicy),
		Status:         status,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&part).Error; err != nil {
			return err
		}
		images := partImagesFromRequest(part.ID, req.Images)
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		fitments := partFitmentsFromRequest(part.ID, req.Fitments)
		if len(fitments) > 0 {
			if err := tx.Create(&fitments).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		serverError(c, err)
		return
	}

	created(c, gin.H{"id": part.ID, "status": part.Status})
}

func (h *PartHandler) Update(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	partID, okID := pathID(c, "partId")
	if !okID {
		return
	}

	var req partRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if err := validatePartRequest(req); err != nil {
		badRequest(c, err.Error())
		return
	}

	status, err := contentStatus(h.db, req.Title, req.Brand, req.Model, req.Description)
	if err != nil {
		serverError(c, err)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Part{}).
			Where("id = ? AND user_id = ? AND status <> ?", partID, userID, statusDeleted).
			Updates(map[string]any{
				"type":            req.Type,
				"category_id":     req.CategoryID,
				"title":           normalizeText(req.Title, 120),
				"brand":           normalizeText(req.Brand, 64),
				"model":           normalizeText(req.Model, 128),
				"condition_level": req.ConditionLevel,
				"price":           req.Price,
				"city_code":       normalizeText(req.CityCode, 32),
				"city_name":       normalizeText(req.CityName, 64),
				"description":     strings.TrimSpace(req.Description),
				"contact_policy":  contactPolicyOrDefault(req.ContactPolicy),
				"status":          status,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Where("part_id = ?", partID).Delete(&model.PartImage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("part_id = ?", partID).Delete(&model.PartFitment{}).Error; err != nil {
			return err
		}
		images := partImagesFromRequest(partID, req.Images)
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		fitments := partFitmentsFromRequest(partID, req.Fitments)
		if len(fitments) > 0 {
			if err := tx.Create(&fitments).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{"id": partID, "status": status})
}

func (h *PartHandler) Delete(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	partID, okID := pathID(c, "partId")
	if !okID {
		return
	}

	result := h.db.Model(&model.Part{}).
		Where("id = ? AND user_id = ? AND status <> ?", partID, userID, statusDeleted).
		Update("status", statusDeleted)
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "part")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PartHandler) Favorite(c *gin.Context) {
	h.togglePartFavorite(c, true)
}

func (h *PartHandler) Unfavorite(c *gin.Context) {
	h.togglePartFavorite(c, false)
}

func (h *PartHandler) togglePartFavorite(c *gin.Context, add bool) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	partID, okID := pathID(c, "partId")
	if !okID {
		return
	}

	var count int64
	if err := h.db.Model(&model.Part{}).Where("id = ? AND status = ?", partID, statusVisible).Count(&count).Error; err != nil {
		serverError(c, err)
		return
	}
	if count == 0 {
		notFound(c, "part")
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if add {
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.PartFavorite{PartID: partID, UserID: userID})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return nil
			}
			return tx.Model(&model.Part{}).Where("id = ?", partID).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
		}

		result := tx.Where("part_id = ? AND user_id = ?", partID, userID).Delete(&model.PartFavorite{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}
		return tx.Model(&model.Part{}).
			Where("id = ? AND favorite_count > 0", partID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
	})
	if err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"id": partID})
}

func (h *PartHandler) ListCategories(c *gin.Context) {
	var categories []model.PartCategory
	if err := h.db.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, categories)
}

func validatePartRequest(req partRequest) error {
	if req.Type < 1 || req.Type > 2 {
		return errors.New("type must be 1 or 2")
	}
	if req.CategoryID <= 0 {
		return errors.New("categoryId is required")
	}
	if normalizeText(req.Title, 120) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(req.Description) == "" {
		return errors.New("description is required")
	}
	if req.ConditionLevel < 0 || req.ConditionLevel > 4 {
		return errors.New("conditionLevel must be between 0 and 4")
	}
	if req.Price != nil && *req.Price < 0 {
		return errors.New("price cannot be negative")
	}
	if len(req.Images) > 9 {
		return errors.New("images cannot exceed 9")
	}
	if len(req.Fitments) > 20 {
		return errors.New("fitments cannot exceed 20")
	}
	return nil
}

func contactPolicyOrDefault(value int8) int8 {
	if value <= 0 {
		return 1
	}
	return value
}

func partImagesFromRequest(partID int64, images []string) []model.PartImage {
	result := make([]model.PartImage, 0, len(images))
	for i, imageURL := range images {
		imageURL = normalizeText(imageURL, 512)
		if imageURL == "" {
			continue
		}
		result = append(result, model.PartImage{
			PartID:    partID,
			ImageURL:  imageURL,
			SortOrder: i,
		})
	}
	return result
}

func partFitmentsFromRequest(partID int64, fitments []fitmentRequest) []model.PartFitment {
	result := make([]model.PartFitment, 0, len(fitments))
	seen := make(map[int64]struct{}, len(fitments))
	for _, fitment := range fitments {
		if fitment.VehicleModelID <= 0 {
			continue
		}
		if _, exists := seen[fitment.VehicleModelID]; exists {
			continue
		}
		seen[fitment.VehicleModelID] = struct{}{}
		result = append(result, model.PartFitment{
			PartID:         partID,
			VehicleModelID: fitment.VehicleModelID,
			Note:           normalizeText(fitment.Note, 255),
		})
	}
	return result
}

func queryFloat(c *gin.Context, name string) (float64, bool) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return 0, false
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

package handler

import (
	"net/http"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(deps Deps) *UserHandler {
	return &UserHandler{db: deps.DB}
}

type updateMeRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

func (h *UserHandler) Me(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}

	updates := map[string]any{
		"nickname":   normalizeText(req.Nickname, 64),
		"avatar_url": normalizeText(req.AvatarURL, 512),
	}
	if err := h.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		serverError(c, err)
		return
	}

	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, user)
}

func (h *UserHandler) ListGarages(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var garages []model.UserGarage
	if err := h.db.Where("user_id = ?", userID).Order("is_primary DESC, id DESC").Find(&garages).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, garages)
}

type garageRequest struct {
	VehicleModelID int64  `json:"vehicleModelId"`
	Year           *int16 `json:"year"`
	Nickname       string `json:"nickname"`
	Description    string `json:"description"`
	IsPrimary      bool   `json:"isPrimary"`
}

func (h *UserHandler) CreateGarage(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var req garageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if req.VehicleModelID <= 0 {
		badRequest(c, "vehicleModelId is required")
		return
	}

	var count int64
	if err := h.db.Model(&model.VehicleModel{}).Where("id = ?", req.VehicleModelID).Count(&count).Error; err != nil {
		serverError(c, err)
		return
	}
	if count == 0 {
		notFound(c, "vehicle model")
		return
	}

	garage := model.UserGarage{
		UserID:         userID,
		VehicleModelID: req.VehicleModelID,
		Year:           req.Year,
		Nickname:       normalizeText(req.Nickname, 64),
		Description:    normalizeText(req.Description, 512),
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		shouldPrimary := req.IsPrimary
		if !shouldPrimary {
			var existing int64
			if err := tx.Model(&model.UserGarage{}).Where("user_id = ?", userID).Count(&existing).Error; err != nil {
				return err
			}
			shouldPrimary = existing == 0
		}
		if shouldPrimary {
			if err := tx.Model(&model.UserGarage{}).Where("user_id = ?", userID).Update("is_primary", 0).Error; err != nil {
				return err
			}
			garage.IsPrimary = 1
		}
		return tx.Create(&garage).Error
	})
	if err != nil {
		serverError(c, err)
		return
	}

	created(c, garage)
}

func (h *UserHandler) UpdateGarage(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	garageID, okID := pathID(c, "garageId")
	if !okID {
		return
	}

	var req garageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if req.VehicleModelID <= 0 {
		badRequest(c, "vehicleModelId is required")
		return
	}

	var garage model.UserGarage
	if err := h.db.Where("id = ? AND user_id = ?", garageID, userID).First(&garage).Error; err != nil {
		serverError(c, err)
		return
	}

	updates := map[string]any{
		"vehicle_model_id": req.VehicleModelID,
		"year":             req.Year,
		"nickname":         normalizeText(req.Nickname, 64),
		"description":      normalizeText(req.Description, 512),
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if req.IsPrimary {
			if err := tx.Model(&model.UserGarage{}).Where("user_id = ?", userID).Update("is_primary", 0).Error; err != nil {
				return err
			}
			updates["is_primary"] = 1
		}
		return tx.Model(&model.UserGarage{}).Where("id = ? AND user_id = ?", garageID, userID).Updates(updates).Error
	})
	if err != nil {
		serverError(c, err)
		return
	}

	if err := h.db.First(&garage, garageID).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, garage)
}

func (h *UserHandler) DeleteGarage(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	garageID, okID := pathID(c, "garageId")
	if !okID {
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", garageID, userID).Delete(&model.UserGarage{})
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "garage")
		return
	}
	c.Status(http.StatusNoContent)
}

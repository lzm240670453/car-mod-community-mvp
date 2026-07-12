package handler

import (
	"strings"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	targetPost     int8 = 1
	targetPart     int8 = 3
	targetUser     int8 = 4
	targetReport   int8 = 5
	targetVehicle  int8 = 6
	targetCategory int8 = 7
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(deps Deps) *AdminHandler {
	return &AdminHandler{db: deps.DB}
}

type adminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	username := strings.TrimSpace(req.Username)
	if username == "" {
		badRequest(c, "username is required")
		return
	}

	var admin model.AdminUser
	if err := h.db.Where("username = ? AND status = ?", username, 1).First(&admin).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{
		"admin": admin,
		"token": "dev-admin-token",
	})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, pageSize := pagination(c)
	query := h.db.Model(&model.User{})

	if status, ok := queryInt8(c, "status"); ok {
		query = query.Where("status = ?", status)
	}
	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("nickname LIKE ? OR phone LIKE ? OR openid LIKE ?", pattern, pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var users []model.User
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"items": users, "page": page, "pageSize": pageSize, "total": total})
}

type statusRequest struct {
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID, okID := pathID(c, "userId")
	if !okID {
		return
	}

	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if req.Status != 1 && req.Status != 2 {
		badRequest(c, "status must be 1 or 2")
		return
	}

	result := h.db.Model(&model.User{}).Where("id = ?", userID).Update("status", req.Status)
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "user")
		return
	}
	_ = h.audit(c, targetUser, userID, "update_status", req.Remark)
	ok(c, gin.H{"id": userID, "status": req.Status})
}

func (h *AdminHandler) ListPostsForReview(c *gin.Context) {
	h.listPostsByStatuses(c, []int8{statusPendingReview, statusVisible, statusHidden})
}

func (h *AdminHandler) ListPendingReviewPosts(c *gin.Context) {
	h.listPostsByStatuses(c, []int8{statusPendingReview})
}

func (h *AdminHandler) listPostsByStatuses(c *gin.Context, statuses []int8) {
	page, pageSize := pagination(c)
	query := h.db.Model(&model.Post{}).Where("status IN ?", statuses)

	if postType, ok := queryInt8(c, "type"); ok {
		query = query.Where("type = ?", postType)
	}
	if userID, ok := queryInt64(c, "userId"); ok {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var posts []model.Post
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"items": posts, "page": page, "pageSize": pageSize, "total": total})
}

func (h *AdminHandler) ApprovePost(c *gin.Context) {
	h.updatePostStatus(c, statusVisible, "approve")
}

func (h *AdminHandler) HidePost(c *gin.Context) {
	h.updatePostStatus(c, statusHidden, "hide")
}

func (h *AdminHandler) RestorePost(c *gin.Context) {
	h.updatePostStatus(c, statusVisible, "restore")
}

func (h *AdminHandler) updatePostStatus(c *gin.Context, status int8, action string) {
	postID, okID := pathID(c, "postId")
	if !okID {
		return
	}

	var req statusRequest
	_ = c.ShouldBindJSON(&req)

	result := h.db.Model(&model.Post{}).Where("id = ?", postID).Update("status", status)
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "post")
		return
	}
	_ = h.audit(c, targetPost, postID, action, req.Remark)
	ok(c, gin.H{"id": postID, "status": status})
}

func (h *AdminHandler) ListPartsForReview(c *gin.Context) {
	h.listPartsByStatuses(c, []int8{statusPendingReview, statusVisible, statusHidden})
}

func (h *AdminHandler) ListPendingReviewParts(c *gin.Context) {
	h.listPartsByStatuses(c, []int8{statusPendingReview})
}

func (h *AdminHandler) listPartsByStatuses(c *gin.Context, statuses []int8) {
	page, pageSize := pagination(c)
	query := h.db.Model(&model.Part{}).Where("status IN ?", statuses)

	if partType, ok := queryInt8(c, "type"); ok {
		query = query.Where("type = ?", partType)
	}
	if categoryID, ok := queryInt64(c, "categoryId"); ok {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var parts []model.Part
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&parts).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"items": parts, "page": page, "pageSize": pageSize, "total": total})
}

func (h *AdminHandler) ApprovePart(c *gin.Context) {
	h.updatePartStatus(c, statusVisible, "approve")
}

func (h *AdminHandler) HidePart(c *gin.Context) {
	h.updatePartStatus(c, statusHidden, "hide")
}

func (h *AdminHandler) RestorePart(c *gin.Context) {
	h.updatePartStatus(c, statusVisible, "restore")
}

func (h *AdminHandler) updatePartStatus(c *gin.Context, status int8, action string) {
	partID, okID := pathID(c, "partId")
	if !okID {
		return
	}

	var req statusRequest
	_ = c.ShouldBindJSON(&req)

	result := h.db.Model(&model.Part{}).Where("id = ?", partID).Update("status", status)
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "part")
		return
	}
	_ = h.audit(c, targetPart, partID, action, req.Remark)
	ok(c, gin.H{"id": partID, "status": status})
}

func (h *AdminHandler) ListReports(c *gin.Context) {
	page, pageSize := pagination(c)
	query := h.db.Model(&model.Report{})

	if status, ok := queryInt8(c, "status"); ok {
		query = query.Where("status = ?", status)
	}
	if targetType, ok := queryInt8(c, "targetType"); ok {
		query = query.Where("target_type = ?", targetType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var reports []model.Report
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reports).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"items": reports, "page": page, "pageSize": pageSize, "total": total})
}

type processReportRequest struct {
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

func (h *AdminHandler) ProcessReport(c *gin.Context) {
	reportID, okID := pathID(c, "reportId")
	if !okID {
		return
	}

	var req processReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if req.Status != 1 && req.Status != 2 {
		badRequest(c, "status must be 1 or 2")
		return
	}

	result := h.db.Model(&model.Report{}).Where("id = ?", reportID).Update("status", req.Status)
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "report")
		return
	}
	_ = h.audit(c, targetReport, reportID, "process", req.Remark)
	ok(c, gin.H{"id": reportID, "status": req.Status})
}

type vehicleAdminItem struct {
	ModelID    int64  `json:"modelId"`
	ModelName  string `json:"modelName"`
	Year       *int16 `json:"year,omitempty"`
	Generation string `json:"generation"`
	Engine     string `json:"engine"`
	SeriesID   int64  `json:"seriesId"`
	SeriesName string `json:"seriesName"`
	BrandID    int64  `json:"brandId"`
	BrandName  string `json:"brandName"`
}

func (h *AdminHandler) ListVehicles(c *gin.Context) {
	page, pageSize := pagination(c)
	query := h.db.Table("vehicle_models AS vm").
		Select("vm.id AS model_id, vm.name AS model_name, vm.year, vm.generation, vm.engine, vs.id AS series_id, vs.name AS series_name, vb.id AS brand_id, vb.name AS brand_name").
		Joins("JOIN vehicle_series AS vs ON vs.id = vm.series_id").
		Joins("JOIN vehicle_brands AS vb ON vb.id = vs.brand_id")

	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("vm.name LIKE ? OR vs.name LIKE ? OR vb.name LIKE ?", pattern, pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var items []vehicleAdminItem
	if err := query.Order("vb.sort_order ASC, vb.id ASC, vs.id ASC, vm.year DESC, vm.id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&items).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"items": items, "page": page, "pageSize": pageSize, "total": total})
}

type brandRequest struct {
	Name      string `json:"name"`
	Initial   string `json:"initial"`
	LogoURL   string `json:"logoUrl"`
	SortOrder int    `json:"sortOrder"`
}

func (h *AdminHandler) CreateBrand(c *gin.Context) {
	var req brandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		badRequest(c, "name is required")
		return
	}

	brand := model.VehicleBrand{
		Name:      normalizeText(req.Name, 64),
		Initial:   normalizeText(req.Initial, 1),
		LogoURL:   normalizeText(req.LogoURL, 512),
		SortOrder: req.SortOrder,
	}
	if err := h.db.Create(&brand).Error; err != nil {
		serverError(c, err)
		return
	}
	_ = h.audit(c, targetVehicle, brand.ID, "create_vehicle_brand", brand.Name)
	created(c, brand)
}

type seriesRequest struct {
	BrandID int64  `json:"brandId"`
	Name    string `json:"name"`
}

func (h *AdminHandler) CreateSeries(c *gin.Context) {
	var req seriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if req.BrandID <= 0 || strings.TrimSpace(req.Name) == "" {
		badRequest(c, "brandId and name are required")
		return
	}

	series := model.VehicleSeries{BrandID: req.BrandID, Name: normalizeText(req.Name, 64)}
	if err := h.db.Create(&series).Error; err != nil {
		serverError(c, err)
		return
	}
	_ = h.audit(c, targetVehicle, series.ID, "create_vehicle_series", series.Name)
	created(c, series)
}

type modelRequest struct {
	SeriesID   int64  `json:"seriesId"`
	Name       string `json:"name"`
	Year       *int16 `json:"year"`
	Generation string `json:"generation"`
	Engine     string `json:"engine"`
}

func (h *AdminHandler) CreateModel(c *gin.Context) {
	var req modelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if req.SeriesID <= 0 || strings.TrimSpace(req.Name) == "" {
		badRequest(c, "seriesId and name are required")
		return
	}

	vehicleModel := model.VehicleModel{
		SeriesID:   req.SeriesID,
		Name:       normalizeText(req.Name, 128),
		Year:       req.Year,
		Generation: normalizeText(req.Generation, 64),
		Engine:     normalizeText(req.Engine, 64),
	}
	if err := h.db.Create(&vehicleModel).Error; err != nil {
		serverError(c, err)
		return
	}
	_ = h.audit(c, targetVehicle, vehicleModel.ID, "create_vehicle_model", vehicleModel.Name)
	created(c, vehicleModel)
}

func (h *AdminHandler) ListCategories(c *gin.Context) {
	var categories []model.PartCategory
	if err := h.db.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, categories)
}

type categoryRequest struct {
	ParentID  *int64 `json:"parentId"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
}

func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		badRequest(c, "name is required")
		return
	}

	category := model.PartCategory{
		ParentID:  req.ParentID,
		Name:      normalizeText(req.Name, 64),
		SortOrder: req.SortOrder,
	}
	if err := h.db.Create(&category).Error; err != nil {
		serverError(c, err)
		return
	}
	_ = h.audit(c, targetCategory, category.ID, "create_category", category.Name)
	created(c, category)
}

func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	categoryID, okID := pathID(c, "categoryId")
	if !okID {
		return
	}

	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		badRequest(c, "name is required")
		return
	}

	result := h.db.Model(&model.PartCategory{}).Where("id = ?", categoryID).Updates(map[string]any{
		"parent_id":  req.ParentID,
		"name":       normalizeText(req.Name, 64),
		"sort_order": req.SortOrder,
	})
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "category")
		return
	}

	var category model.PartCategory
	if err := h.db.First(&category, categoryID).Error; err != nil {
		serverError(c, err)
		return
	}
	_ = h.audit(c, targetCategory, category.ID, "update_category", category.Name)
	ok(c, category)
}

func (h *AdminHandler) audit(c *gin.Context, targetType int8, targetID int64, action string, remark string) error {
	adminID := currentAdminID(c)
	if adminID == 0 {
		adminID = 1
	}
	return h.db.Create(&model.AuditLog{
		AdminID:    adminID,
		TargetType: targetType,
		TargetID:   targetID,
		Action:     normalizeText(action, 64),
		Remark:     normalizeText(remark, 500),
	}).Error
}

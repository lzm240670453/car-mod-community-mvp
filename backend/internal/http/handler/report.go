package handler

import (
	"errors"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(deps Deps) *ReportHandler {
	return &ReportHandler{db: deps.DB}
}

type reportRequest struct {
	TargetType int8   `json:"targetType"`
	TargetID   int64  `json:"targetId"`
	ReasonType int8   `json:"reasonType"`
	ReasonText string `json:"reasonText"`
}

func (h *ReportHandler) Create(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var req reportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if err := validateReportRequest(req); err != nil {
		badRequest(c, err.Error())
		return
	}

	report := model.Report{
		ReporterID: userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		ReasonType: req.ReasonType,
		ReasonText: normalizeText(req.ReasonText, 500),
		Status:     0,
	}
	if err := h.db.Create(&report).Error; err != nil {
		serverError(c, err)
		return
	}
	created(c, gin.H{"id": report.ID, "status": report.Status})
}

func (h *ReportHandler) My(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	page, pageSize := pagination(c)
	query := h.db.Model(&model.Report{}).Where("reporter_id = ?", userID)

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

	ok(c, gin.H{
		"items":    reports,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func validateReportRequest(req reportRequest) error {
	if req.TargetType < 1 || req.TargetType > 4 {
		return errors.New("targetType must be between 1 and 4")
	}
	if req.TargetID <= 0 {
		return errors.New("targetId is required")
	}
	if req.ReasonType <= 0 {
		return errors.New("reasonType is required")
	}
	if normalizeText(req.ReasonText, 500) == "" {
		return errors.New("reasonText is required")
	}
	return nil
}

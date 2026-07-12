package handler

import (
	"net/http"
	"time"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageHandler struct {
	db *gorm.DB
}

func NewMessageHandler(deps Deps) *MessageHandler {
	return &MessageHandler{db: deps.DB}
}

func (h *MessageHandler) List(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	page, pageSize := pagination(c)
	query := h.db.Model(&model.SiteMessage{}).Where("recipient_id = ?", userID)

	if messageType, ok := queryInt8(c, "type"); ok {
		query = query.Where("type = ?", messageType)
	}
	if c.Query("unread") == "1" {
		query = query.Where("read_at IS NULL")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var messages []model.SiteMessage
	if err := query.Order("read_at IS NULL DESC, created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{
		"items":    messages,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (h *MessageHandler) UnreadCount(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var total int64
	if err := h.db.Model(&model.SiteMessage{}).
		Where("recipient_id = ? AND read_at IS NULL", userID).
		Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"count": total})
}

func (h *MessageHandler) MarkRead(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	messageID, okID := pathID(c, "messageId")
	if !okID {
		return
	}

	var count int64
	if err := h.db.Model(&model.SiteMessage{}).
		Where("id = ? AND recipient_id = ?", messageID, userID).
		Count(&count).Error; err != nil {
		serverError(c, err)
		return
	}
	if count == 0 {
		notFound(c, "message")
		return
	}

	if err := h.db.Model(&model.SiteMessage{}).
		Where("id = ? AND recipient_id = ? AND read_at IS NULL", messageID, userID).
		Update("read_at", time.Now()).Error; err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *MessageHandler) MarkAllRead(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	if err := h.db.Model(&model.SiteMessage{}).
		Where("recipient_id = ? AND read_at IS NULL", userID).
		Update("read_at", time.Now()).Error; err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

package handler

import (
	"errors"
	"net/http"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntentHandler struct {
	db *gorm.DB
}

func NewIntentHandler(deps Deps) *IntentHandler {
	return &IntentHandler{db: deps.DB}
}

type intentRequest struct {
	Message string `json:"message"`
}

func (h *IntentHandler) Create(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	partID, okID := pathID(c, "partId")
	if !okID {
		return
	}

	var req intentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}

	var part model.Part
	if err := h.db.Where("id = ? AND status = ?", partID, statusVisible).First(&part).Error; err != nil {
		serverError(c, err)
		return
	}
	if part.UserID == userID {
		badRequest(c, "cannot create intent for your own part")
		return
	}

	var intent model.TradeIntent
	err := h.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("part_id = ? AND buyer_id = ?", partID, userID).First(&intent).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil {
			intent.Message = normalizeText(req.Message, 500)
			intent.Status = 1
			if err := tx.Model(&model.TradeIntent{}).
				Where("id = ?", intent.ID).
				Updates(map[string]any{
					"message": intent.Message,
					"status":  1,
				}).Error; err != nil {
				return err
			}
			actorID := userID
			content := intent.Message
			if content == "" {
				content = "对方更新了这条二手件交易意向"
			}
			return createSiteMessage(tx, part.UserID, &actorID, messageTypeTrade, "二手件意向已更新", content, messageTargetPart, partID)
		}

		intent = model.TradeIntent{
			PartID:   partID,
			BuyerID:  userID,
			SellerID: part.UserID,
			Message:  normalizeText(req.Message, 500),
			Status:   1,
		}
		if err := tx.Create(&intent).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Part{}).Where("id = ?", partID).UpdateColumn("intent_count", gorm.Expr("intent_count + ?", 1)).Error; err != nil {
			return err
		}
		actorID := userID
		content := intent.Message
		if content == "" {
			content = "有车友对你的二手件发送了交易意向"
		}
		return createSiteMessage(tx, part.UserID, &actorID, messageTypeTrade, "新的二手件意向", content, messageTargetPart, partID)
	})
	if err != nil {
		serverError(c, err)
		return
	}

	created(c, gin.H{"intentId": intent.ID, "status": intent.Status})
}

func (h *IntentHandler) List(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	page, pageSize := pagination(c)
	query := h.db.Model(&model.TradeIntent{}).Where("buyer_id = ? OR seller_id = ?", userID, userID)

	if status, ok := queryInt8(c, "status"); ok {
		query = query.Where("status = ?", status)
	}
	if role := c.Query("role"); role == "buyer" {
		query = query.Where("buyer_id = ?", userID)
	} else if role == "seller" {
		query = query.Where("seller_id = ?", userID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var intents []model.TradeIntent
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&intents).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{
		"items":    intents,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (h *IntentHandler) Get(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	intentID, okID := pathID(c, "intentId")
	if !okID {
		return
	}

	var intent model.TradeIntent
	if err := h.db.Where("id = ? AND (buyer_id = ? OR seller_id = ?)", intentID, userID, userID).First(&intent).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, intent)
}

func (h *IntentHandler) Close(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	intentID, okID := pathID(c, "intentId")
	if !okID {
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var intent model.TradeIntent
		if err := tx.Where("id = ? AND (buyer_id = ? OR seller_id = ?) AND status = ?", intentID, userID, userID, 1).First(&intent).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.TradeIntent{}).Where("id = ?", intentID).Update("status", 3).Error; err != nil {
			return err
		}

		recipientID := intent.SellerID
		if userID == intent.SellerID {
			recipientID = intent.BuyerID
		}
		actorID := userID
		return createSiteMessage(tx, recipientID, &actorID, messageTypeTrade, "交易意向已关闭", "对方关闭了这条二手件交易意向", messageTargetPart, intent.PartID)
	})
	if err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

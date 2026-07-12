package handler

import (
	"retrofit/backend/internal/model"

	"gorm.io/gorm"
)

const (
	messageTypeSystem      int8 = 1
	messageTypeTrade       int8 = 2
	messageTypeInteraction int8 = 3

	messageTargetNone    int8 = 0
	messageTargetPost    int8 = 1
	messageTargetComment int8 = 2
	messageTargetPart    int8 = 3
	messageTargetUser    int8 = 4
)

func createSiteMessage(tx *gorm.DB, recipientID int64, actorID *int64, messageType int8, title string, content string, targetType int8, targetID int64) error {
	if recipientID <= 0 {
		return nil
	}
	if actorID != nil && *actorID == recipientID {
		return nil
	}
	if targetID <= 0 {
		targetType = messageTargetNone
		targetID = 0
	}

	message := model.SiteMessage{
		RecipientID: recipientID,
		ActorID:     actorID,
		Type:        messageType,
		Title:       normalizeText(title, 120),
		Content:     normalizeText(content, 500),
		TargetType:  targetType,
		TargetID:    targetID,
	}
	if message.Title == "" {
		message.Title = "新的站内信"
	}
	return tx.Create(&message).Error
}

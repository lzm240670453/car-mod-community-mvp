package model

import "time"

type SiteMessage struct {
	ID          int64      `gorm:"column:id;primaryKey" json:"id"`
	RecipientID int64      `gorm:"column:recipient_id" json:"recipientId"`
	ActorID     *int64     `gorm:"column:actor_id" json:"actorId,omitempty"`
	Type        int8       `gorm:"column:type" json:"type"`
	Title       string     `gorm:"column:title" json:"title"`
	Content     string     `gorm:"column:content" json:"content"`
	TargetType  int8       `gorm:"column:target_type" json:"targetType"`
	TargetID    int64      `gorm:"column:target_id" json:"targetId"`
	ReadAt      *time.Time `gorm:"column:read_at" json:"readAt,omitempty"`
	Timestamp
}

func (SiteMessage) TableName() string {
	return "site_messages"
}

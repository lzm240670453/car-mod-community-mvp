package model

import "time"

type Report struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"id"`
	ReporterID int64  `gorm:"column:reporter_id" json:"reporterId"`
	TargetType int8   `gorm:"column:target_type" json:"targetType"`
	TargetID   int64  `gorm:"column:target_id" json:"targetId"`
	ReasonType int8   `gorm:"column:reason_type" json:"reasonType"`
	ReasonText string `gorm:"column:reason_text" json:"reasonText"`
	Status     int8   `gorm:"column:status" json:"status"`
	Timestamp
}

func (Report) TableName() string {
	return "reports"
}

type AuditLog struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`
	AdminID    int64     `gorm:"column:admin_id" json:"adminId"`
	TargetType int8      `gorm:"column:target_type" json:"targetType"`
	TargetID   int64     `gorm:"column:target_id" json:"targetId"`
	Action     string    `gorm:"column:action" json:"action"`
	Remark     string    `gorm:"column:remark" json:"remark"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

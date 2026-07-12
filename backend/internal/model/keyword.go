package model

type HighRiskKeyword struct {
	ID      int64  `gorm:"column:id;primaryKey" json:"id"`
	Keyword string `gorm:"column:keyword" json:"keyword"`
	Action  int8   `gorm:"column:action" json:"action"`
	Enabled int8   `gorm:"column:enabled" json:"enabled"`
	Timestamp
}

func (HighRiskKeyword) TableName() string {
	return "high_risk_keywords"
}

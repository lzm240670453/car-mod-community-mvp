package model

type AdminUser struct {
	ID           int64  `gorm:"column:id;primaryKey" json:"id"`
	Username     string `gorm:"column:username" json:"username"`
	PasswordHash string `gorm:"column:password_hash" json:"-"`
	Role         string `gorm:"column:role" json:"role"`
	Status       int8   `gorm:"column:status" json:"status"`
	Timestamp
}

func (AdminUser) TableName() string {
	return "admin_users"
}

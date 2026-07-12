package model

import "time"

type User struct {
	ID           int64      `gorm:"column:id;primaryKey" json:"id"`
	OpenID       string     `gorm:"column:openid" json:"openid"`
	UnionID      *string    `gorm:"column:unionid" json:"unionid,omitempty"`
	Nickname     string     `gorm:"column:nickname" json:"nickname"`
	AvatarURL    string     `gorm:"column:avatar_url" json:"avatarUrl"`
	Phone        string     `gorm:"column:phone" json:"phone"`
	PhoneBoundAt *time.Time `gorm:"column:phone_bound_at" json:"phoneBoundAt,omitempty"`
	Status       int8       `gorm:"column:status" json:"status"`
	Timestamp
}

func (User) TableName() string {
	return "users"
}

type UserGarage struct {
	ID             int64  `gorm:"column:id;primaryKey" json:"id"`
	UserID         int64  `gorm:"column:user_id" json:"userId"`
	VehicleModelID int64  `gorm:"column:vehicle_model_id" json:"vehicleModelId"`
	Year           *int16 `gorm:"column:year" json:"year,omitempty"`
	Nickname       string `gorm:"column:nickname" json:"nickname"`
	Description    string `gorm:"column:description" json:"description"`
	IsPrimary      int8   `gorm:"column:is_primary" json:"isPrimary"`
	Timestamp
}

func (UserGarage) TableName() string {
	return "user_garages"
}

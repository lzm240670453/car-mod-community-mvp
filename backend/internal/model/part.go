package model

import (
	"time"
)

type PartCategory struct {
	ID        int64  `gorm:"column:id;primaryKey" json:"id"`
	ParentID  *int64 `gorm:"column:parent_id" json:"parentId,omitempty"`
	Name      string `gorm:"column:name" json:"name"`
	SortOrder int    `gorm:"column:sort_order" json:"sortOrder"`
	Timestamp
}

func (PartCategory) TableName() string {
	return "part_categories"
}

type Part struct {
	ID             int64    `gorm:"column:id;primaryKey" json:"id"`
	UserID         int64    `gorm:"column:user_id" json:"userId"`
	Type           int8     `gorm:"column:type" json:"type"`
	CategoryID     int64    `gorm:"column:category_id" json:"categoryId"`
	Title          string   `gorm:"column:title" json:"title"`
	Brand          string   `gorm:"column:brand" json:"brand"`
	Model          string   `gorm:"column:model" json:"model"`
	ConditionLevel int8     `gorm:"column:condition_level" json:"conditionLevel"`
	Price          *float64 `gorm:"column:price" json:"price,omitempty"`
	CityCode       string   `gorm:"column:city_code" json:"cityCode"`
	CityName       string   `gorm:"column:city_name" json:"cityName"`
	Description    string   `gorm:"column:description" json:"description"`
	ContactPolicy  int8     `gorm:"column:contact_policy" json:"contactPolicy"`
	Status         int8     `gorm:"column:status" json:"status"`
	ViewCount      int      `gorm:"column:view_count" json:"viewCount"`
	FavoriteCount  int      `gorm:"column:favorite_count" json:"favoriteCount"`
	IntentCount    int      `gorm:"column:intent_count" json:"intentCount"`
	Timestamp
}

func (Part) TableName() string {
	return "parts"
}

type PartImage struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	PartID    int64     `gorm:"column:part_id" json:"partId"`
	ImageURL  string    `gorm:"column:image_url" json:"imageUrl"`
	SortOrder int       `gorm:"column:sort_order" json:"sortOrder"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (PartImage) TableName() string {
	return "part_images"
}

type PartFitment struct {
	ID             int64     `gorm:"column:id;primaryKey" json:"id"`
	PartID         int64     `gorm:"column:part_id" json:"partId"`
	VehicleModelID int64     `gorm:"column:vehicle_model_id" json:"vehicleModelId"`
	Note           string    `gorm:"column:note" json:"note"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (PartFitment) TableName() string {
	return "part_fitments"
}

type PartFavorite struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	PartID    int64     `gorm:"column:part_id" json:"partId"`
	UserID    int64     `gorm:"column:user_id" json:"userId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (PartFavorite) TableName() string {
	return "part_favorites"
}

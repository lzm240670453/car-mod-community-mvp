package model

type VehicleBrand struct {
	ID        int64  `gorm:"column:id;primaryKey" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	Initial   string `gorm:"column:initial" json:"initial"`
	LogoURL   string `gorm:"column:logo_url" json:"logoUrl"`
	SortOrder int    `gorm:"column:sort_order" json:"sortOrder"`
	Timestamp
}

func (VehicleBrand) TableName() string {
	return "vehicle_brands"
}

type VehicleSeries struct {
	ID      int64  `gorm:"column:id;primaryKey" json:"id"`
	BrandID int64  `gorm:"column:brand_id" json:"brandId"`
	Name    string `gorm:"column:name" json:"name"`
	Timestamp
}

func (VehicleSeries) TableName() string {
	return "vehicle_series"
}

type VehicleModel struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"id"`
	SeriesID   int64  `gorm:"column:series_id" json:"seriesId"`
	Name       string `gorm:"column:name" json:"name"`
	Year       *int16 `gorm:"column:year" json:"year,omitempty"`
	Generation string `gorm:"column:generation" json:"generation"`
	Engine     string `gorm:"column:engine" json:"engine"`
	Timestamp
}

func (VehicleModel) TableName() string {
	return "vehicle_models"
}

package handler

import (
	"net/http"
	"strings"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VehicleHandler struct {
	db *gorm.DB
}

func NewVehicleHandler(deps Deps) *VehicleHandler {
	return &VehicleHandler{db: deps.DB}
}

func (h *VehicleHandler) ListBrands(c *gin.Context) {
	var brands []model.VehicleBrand
	if err := h.db.Order("sort_order ASC, initial ASC, id ASC").Find(&brands).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, brands)
}

func (h *VehicleHandler) ListSeries(c *gin.Context) {
	brandID, okID := pathID(c, "brandId")
	if !okID {
		return
	}

	var series []model.VehicleSeries
	if err := h.db.Where("brand_id = ?", brandID).Order("id ASC").Find(&series).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, series)
}

func (h *VehicleHandler) ListModels(c *gin.Context) {
	seriesID, okID := pathID(c, "seriesId")
	if !okID {
		return
	}

	var models []model.VehicleModel
	if err := h.db.Where("series_id = ?", seriesID).Order("year DESC, id ASC").Find(&models).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, models)
}

type vehicleSearchItem struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Year       *int16 `json:"year,omitempty"`
	Generation string `json:"generation"`
	Engine     string `json:"engine"`
	SeriesID   int64  `json:"seriesId"`
	SeriesName string `json:"seriesName"`
	BrandID    int64  `json:"brandId"`
	BrandName  string `json:"brandName"`
}

func (h *VehicleHandler) Search(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("q"))
	if keyword == "" {
		badRequest(c, "q is required")
		return
	}

	limit := 20
	var items []vehicleSearchItem
	pattern := "%" + keyword + "%"
	err := h.db.Table("vehicle_models AS vm").
		Select("vm.id, vm.name, vm.year, vm.generation, vm.engine, vm.series_id, vs.name AS series_name, vb.id AS brand_id, vb.name AS brand_name").
		Joins("JOIN vehicle_series AS vs ON vs.id = vm.series_id").
		Joins("JOIN vehicle_brands AS vb ON vb.id = vs.brand_id").
		Where("vm.name LIKE ? OR vm.generation LIKE ? OR vm.engine LIKE ? OR vs.name LIKE ? OR vb.name LIKE ?", pattern, pattern, pattern, pattern, pattern).
		Order("vb.sort_order ASC, vm.year DESC, vm.id ASC").
		Limit(limit).
		Scan(&items).Error
	if err != nil {
		serverError(c, err)
		return
	}
	if items == nil {
		items = []vehicleSearchItem{}
	}
	c.JSON(http.StatusOK, Response{Data: items})
}

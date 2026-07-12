package handler

import (
	"strings"

	"retrofit/backend/internal/model"

	"gorm.io/gorm"
)

const (
	statusPendingReview int8 = 0
	statusVisible       int8 = 1
	statusHidden        int8 = 2
	statusDeleted       int8 = 3
	statusClosed        int8 = 4
)

func contentStatus(db *gorm.DB, values ...string) (int8, error) {
	var keywords []model.HighRiskKeyword
	if err := db.Where("enabled = ?", 1).Find(&keywords).Error; err != nil {
		return statusVisible, err
	}

	text := strings.ToLower(strings.Join(values, "\n"))
	for _, keyword := range keywords {
		if keyword.Keyword == "" {
			continue
		}
		if strings.Contains(text, strings.ToLower(keyword.Keyword)) {
			return statusPendingReview, nil
		}
	}
	return statusVisible, nil
}

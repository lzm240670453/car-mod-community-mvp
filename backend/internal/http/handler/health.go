package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redis}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	status := gin.H{"service": "ok", "mysql": "ok", "redis": "ok"}

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		status["mysql"] = "error"
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		status["redis"] = "error"
	}

	code := http.StatusOK
	if status["mysql"] == "error" {
		code = http.StatusServiceUnavailable
	}
	c.JSON(code, Response{Data: status})
}

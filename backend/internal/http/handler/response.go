package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Data: data})
}

func created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{Data: data})
}

func fail(c *gin.Context, code int, message string) {
	c.JSON(code, Response{Error: message})
}

func badRequest(c *gin.Context, message string) {
	fail(c, http.StatusBadRequest, message)
}

func unauthorized(c *gin.Context) {
	fail(c, http.StatusUnauthorized, "missing or invalid user")
}

func notFound(c *gin.Context, name string) {
	fail(c, http.StatusNotFound, name+" not found")
}

func serverError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fail(c, http.StatusNotFound, "resource not found")
		return
	}
	fail(c, http.StatusInternalServerError, "internal server error")
}

func placeholder(c *gin.Context, name string) {
	ok(c, gin.H{"handler": name, "status": "not_implemented"})
}

func currentUserID(c *gin.Context) (int64, bool) {
	raw := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if raw == "" {
		raw = strings.TrimSpace(c.Query("userId"))
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func currentAdminID(c *gin.Context) int64 {
	raw := strings.TrimSpace(c.GetHeader("X-Admin-ID"))
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0
	}
	return id
}

func pathID(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		badRequest(c, "invalid "+name)
		return 0, false
	}
	return id, true
}

func queryInt64(c *gin.Context, name string) (int64, bool) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func queryInt8(c *gin.Context, name string) (int8, bool) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 8)
	if err != nil {
		return 0, false
	}
	return int8(id), true
}

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func normalizeText(value string, max int) string {
	value = strings.TrimSpace(value)
	if max > 0 && len([]rune(value)) > max {
		runes := []rune(value)
		return string(runes[:max])
	}
	return value
}

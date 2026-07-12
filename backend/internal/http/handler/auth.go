package handler

import (
	"strings"
	"time"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(deps Deps) *AuthHandler {
	return &AuthHandler{db: deps.DB}
}

type wechatLoginRequest struct {
	Code      string `json:"code"`
	OpenID    string `json:"openid"`
	UnionID   string `json:"unionid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	var req wechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}

	openID := strings.TrimSpace(req.OpenID)
	if openID == "" {
		openID = strings.TrimSpace(req.Code)
	}
	if openID == "" {
		badRequest(c, "openid or code is required")
		return
	}

	var user model.User
	err := h.db.Where("openid = ?", openID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		serverError(c, err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		var unionID *string
		if strings.TrimSpace(req.UnionID) != "" {
			value := normalizeText(req.UnionID, 64)
			unionID = &value
		}
		user = model.User{
			OpenID:    openID,
			UnionID:   unionID,
			Nickname:  normalizeText(req.Nickname, 64),
			AvatarURL: normalizeText(req.AvatarURL, 512),
			Status:    1,
		}
		if err := h.db.Create(&user).Error; err != nil {
			serverError(c, err)
			return
		}
	} else if req.Nickname != "" || req.AvatarURL != "" {
		updates := map[string]any{}
		if req.Nickname != "" {
			updates["nickname"] = normalizeText(req.Nickname, 64)
		}
		if req.AvatarURL != "" {
			updates["avatar_url"] = normalizeText(req.AvatarURL, 512)
		}
		if len(updates) > 0 {
			if err := h.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
				serverError(c, err)
				return
			}
			_ = h.db.First(&user, user.ID).Error
		}
	}

	ok(c, gin.H{
		"user":  user,
		"token": "dev-user-token",
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	ok(c, gin.H{"status": "ok"})
}

type bindPhoneRequest struct {
	Phone string `json:"phone"`
}

func (h *AuthHandler) BindPhone(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var req bindPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	phone := normalizeText(req.Phone, 32)
	if phone == "" {
		badRequest(c, "phone is required")
		return
	}

	now := time.Now()
	result := h.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"phone":          phone,
		"phone_bound_at": &now,
	})
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "user")
		return
	}

	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		serverError(c, err)
		return
	}
	ok(c, user)
}

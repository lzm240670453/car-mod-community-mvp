package handler

import (
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

type uploadSignatureRequest struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Scene       string `json:"scene"`
}

func (h *UploadHandler) Signature(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var req uploadSignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}

	scene := normalizeText(req.Scene, 32)
	if scene == "" {
		scene = "common"
	}
	fileName := sanitizeFileName(req.FileName)
	if fileName == "" {
		fileName = "image.jpg"
	}

	key := strings.Join([]string{
		"uploads",
		scene,
		time.Now().Format("20060102"),
		strings.TrimSpace(c.GetHeader("X-User-ID")) + "-" + time.Now().Format("150405000000") + "-" + fileName,
	}, "/")

	ok(c, gin.H{
		"key":         key,
		"uploadUrl":   "/dev-upload/" + key,
		"publicUrl":   "/static/" + key,
		"method":      "PUT",
		"headers":     gin.H{"Content-Type": normalizeText(req.ContentType, 128)},
		"expiresIn":   600,
		"userId":      userID,
		"storageMode": "dev",
	})
}

type uploadCompleteRequest struct {
	Key       string `json:"key"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	MediaType string `json:"mediaType"`
}

func (h *UploadHandler) Complete(c *gin.Context) {
	if _, okUser := currentUserID(c); !okUser {
		unauthorized(c)
		return
	}

	var req uploadCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	key := normalizeText(req.Key, 512)
	if key == "" {
		badRequest(c, "key is required")
		return
	}

	url := normalizeText(req.URL, 512)
	if url == "" {
		url = "/static/" + key
	}

	ok(c, gin.H{
		"key":       key,
		"url":       url,
		"size":      req.Size,
		"width":     req.Width,
		"height":    req.Height,
		"mediaType": normalizeText(req.MediaType, 64),
	})
}

func sanitizeFileName(fileName string) string {
	fileName = strings.TrimSpace(path.Base(strings.ReplaceAll(fileName, "\\", "/")))
	fileName = strings.ReplaceAll(fileName, " ", "-")
	return normalizeText(fileName, 128)
}

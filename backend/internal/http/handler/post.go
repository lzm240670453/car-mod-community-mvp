package handler

import (
	"errors"
	"net/http"
	"strings"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(deps Deps) *PostHandler {
	return &PostHandler{db: deps.DB}
}

type postRequest struct {
	GarageID       *int64   `json:"garageId"`
	Type           int8     `json:"type"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	VehicleModelID *int64   `json:"vehicleModelId"`
	Images         []string `json:"images"`
}

type postDetail struct {
	model.Post
	Images []model.PostImage `json:"images"`
}

func (h *PostHandler) List(c *gin.Context) {
	page, pageSize := pagination(c)
	query := h.db.Model(&model.Post{}).Where("status = ?", statusVisible)

	if postType, ok := queryInt8(c, "type"); ok {
		query = query.Where("type = ?", postType)
	}
	if vehicleModelID, ok := queryInt64(c, "vehicleModelId"); ok {
		query = query.Where("vehicle_model_id = ?", vehicleModelID)
	}
	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", pattern, pattern)
	}
	if userID, ok := queryInt64(c, "userId"); ok {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var posts []model.Post
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{
		"items":    posts,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (h *PostHandler) Get(c *gin.Context) {
	postID, okID := pathID(c, "postId")
	if !okID {
		return
	}

	var post model.Post
	if err := h.db.Where("id = ? AND status IN ?", postID, []int8{statusVisible, statusPendingReview}).First(&post).Error; err != nil {
		serverError(c, err)
		return
	}

	var images []model.PostImage
	if err := h.db.Where("post_id = ?", postID).Order("sort_order ASC, id ASC").Find(&images).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, postDetail{Post: post, Images: images})
}

func (h *PostHandler) Create(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}

	var req postRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if err := validatePostRequest(req); err != nil {
		badRequest(c, err.Error())
		return
	}

	status, err := contentStatus(h.db, req.Title, req.Content)
	if err != nil {
		serverError(c, err)
		return
	}

	post := model.Post{
		UserID:         userID,
		GarageID:       req.GarageID,
		Type:           req.Type,
		Title:          normalizeText(req.Title, 120),
		Content:        strings.TrimSpace(req.Content),
		VehicleModelID: req.VehicleModelID,
		Status:         status,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		images := postImagesFromRequest(post.ID, req.Images)
		if len(images) > 0 {
			return tx.Create(&images).Error
		}
		return nil
	})
	if err != nil {
		serverError(c, err)
		return
	}

	created(c, gin.H{"id": post.ID, "status": post.Status})
}

func (h *PostHandler) Update(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	postID, okID := pathID(c, "postId")
	if !okID {
		return
	}

	var req postRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	if err := validatePostRequest(req); err != nil {
		badRequest(c, err.Error())
		return
	}

	status, err := contentStatus(h.db, req.Title, req.Content)
	if err != nil {
		serverError(c, err)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Post{}).
			Where("id = ? AND user_id = ? AND status <> ?", postID, userID, statusDeleted).
			Updates(map[string]any{
				"garage_id":        req.GarageID,
				"type":             req.Type,
				"title":            normalizeText(req.Title, 120),
				"content":          strings.TrimSpace(req.Content),
				"vehicle_model_id": req.VehicleModelID,
				"status":           status,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Where("post_id = ?", postID).Delete(&model.PostImage{}).Error; err != nil {
			return err
		}
		images := postImagesFromRequest(postID, req.Images)
		if len(images) > 0 {
			return tx.Create(&images).Error
		}
		return nil
	})
	if err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{"id": postID, "status": status})
}

func (h *PostHandler) Delete(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	postID, okID := pathID(c, "postId")
	if !okID {
		return
	}

	result := h.db.Model(&model.Post{}).
		Where("id = ? AND user_id = ? AND status <> ?", postID, userID, statusDeleted).
		Update("status", statusDeleted)
	if result.Error != nil {
		serverError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		notFound(c, "post")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PostHandler) Like(c *gin.Context) {
	h.togglePostRelation(c, "like", true)
}

func (h *PostHandler) Unlike(c *gin.Context) {
	h.togglePostRelation(c, "like", false)
}

func (h *PostHandler) Favorite(c *gin.Context) {
	h.togglePostRelation(c, "favorite", true)
}

func (h *PostHandler) Unfavorite(c *gin.Context) {
	h.togglePostRelation(c, "favorite", false)
}

func (h *PostHandler) togglePostRelation(c *gin.Context, relation string, add bool) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	postID, okID := pathID(c, "postId")
	if !okID {
		return
	}

	var count int64
	if err := h.db.Model(&model.Post{}).Where("id = ? AND status = ?", postID, statusVisible).Count(&count).Error; err != nil {
		serverError(c, err)
		return
	}
	if count == 0 {
		notFound(c, "post")
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		counter := "like_count"
		if relation == "favorite" {
			counter = "favorite_count"
		}

		if add {
			if relation == "favorite" {
				result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.PostFavorite{PostID: postID, UserID: userID})
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return nil
				}
			} else {
				result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.PostLike{PostID: postID, UserID: userID})
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return nil
				}
			}
			return tx.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn(counter, gorm.Expr(counter+" + ?", 1)).Error
		}

		var result *gorm.DB
		if relation == "favorite" {
			result = tx.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostFavorite{})
		} else {
			result = tx.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostLike{})
		}
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}
		return tx.Model(&model.Post{}).Where("id = ? AND "+counter+" > 0", postID).UpdateColumn(counter, gorm.Expr(counter+" - ?", 1)).Error
	})
	if err != nil {
		serverError(c, err)
		return
	}
	ok(c, gin.H{"id": postID})
}

func (h *PostHandler) ListComments(c *gin.Context) {
	postID, okID := pathID(c, "postId")
	if !okID {
		return
	}

	page, pageSize := pagination(c)
	query := h.db.Model(&model.Comment{}).Where("post_id = ? AND status = ?", postID, statusVisible)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var comments []model.Comment
	if err := query.Order("created_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{
		"items":    comments,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

type commentRequest struct {
	ParentID *int64 `json:"parentId"`
	Content  string `json:"content"`
}

func (h *PostHandler) CreateComment(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	postID, okID := pathID(c, "postId")
	if !okID {
		return
	}

	var req commentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid request body")
		return
	}
	content := normalizeText(req.Content, 1000)
	if content == "" {
		badRequest(c, "content is required")
		return
	}

	var post model.Post
	if err := h.db.Where("id = ? AND status = ?", postID, statusVisible).First(&post).Error; err != nil {
		serverError(c, err)
		return
	}

	status, err := contentStatus(h.db, content)
	if err != nil {
		serverError(c, err)
		return
	}
	comment := model.Comment{
		PostID:   postID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  content,
		Status:   status,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}
		if status == statusVisible {
			if err := tx.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return err
			}
			actorID := userID
			return createSiteMessage(tx, post.UserID, &actorID, messageTypeInteraction, "帖子收到新评论", content, messageTargetPost, postID)
		}
		return nil
	})
	if err != nil {
		serverError(c, err)
		return
	}

	created(c, gin.H{"id": comment.ID, "status": comment.Status})
}

func (h *PostHandler) DeleteComment(c *gin.Context) {
	userID, okUser := currentUserID(c)
	if !okUser {
		unauthorized(c)
		return
	}
	commentID, okID := pathID(c, "commentId")
	if !okID {
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		if err := tx.Where("id = ? AND user_id = ? AND status <> ?", commentID, userID, statusDeleted).First(&comment).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Comment{}).Where("id = ?", commentID).Update("status", statusDeleted).Error; err != nil {
			return err
		}
		if comment.Status == statusVisible {
			return tx.Model(&model.Post{}).
				Where("id = ? AND comment_count > 0", comment.PostID).
				UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		}
		return nil
	})
	if err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func validatePostRequest(req postRequest) error {
	if req.Type < 1 || req.Type > 5 {
		return errors.New("type must be between 1 and 5")
	}
	if normalizeText(req.Title, 120) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("content is required")
	}
	if len(req.Images) > 9 {
		return errors.New("images cannot exceed 9")
	}
	return nil
}

func postImagesFromRequest(postID int64, images []string) []model.PostImage {
	result := make([]model.PostImage, 0, len(images))
	for i, imageURL := range images {
		imageURL = normalizeText(imageURL, 512)
		if imageURL == "" {
			continue
		}
		result = append(result, model.PostImage{
			PostID:    postID,
			ImageURL:  imageURL,
			SortOrder: i,
		})
	}
	return result
}

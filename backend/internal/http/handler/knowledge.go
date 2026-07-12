package handler

import (
	"encoding/json"
	"strings"
	"time"

	"retrofit/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KnowledgeHandler struct {
	db *gorm.DB
}

func NewKnowledgeHandler(deps Deps) *KnowledgeHandler {
	return &KnowledgeHandler{db: deps.DB}
}

type knowledgeCategoryItem struct {
	model.KnowledgeCategory
	ChildCount   int64 `json:"childCount"`
	ArticleCount int64 `json:"articleCount"`
}

type knowledgeArticleItem struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content,omitempty"`
	Tags        []string   `json:"tags"`
	Images      []string   `json:"images"`
	VideoURL    string     `json:"videoUrl"`
	BrandID     *int64     `json:"brandId,omitempty"`
	PriceRange  string     `json:"priceRange"`
	ViewCount   int        `json:"viewCount"`
	LikeCount   int        `json:"likeCount"`
	SortOrder   int        `json:"sortOrder"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type knowledgeCategoryDetail struct {
	model.KnowledgeCategory
	Breadcrumb []model.KnowledgeCategory `json:"breadcrumb"`
	Children   []knowledgeCategoryItem   `json:"children"`
	Articles   []knowledgeArticleItem    `json:"articles"`
}

type knowledgeArticleDetail struct {
	knowledgeArticleItem
	Categories []model.KnowledgeCategory `json:"categories"`
	Brand      *model.KnowledgeBrand     `json:"brand,omitempty"`
}

func (h *KnowledgeHandler) ListCategories(c *gin.Context) {
	query := h.db.Model(&model.KnowledgeCategory{})

	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", pattern, pattern)
	} else if parentID, ok := queryInt64(c, "parentId"); ok {
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	var categories []model.KnowledgeCategory
	if err := query.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		serverError(c, err)
		return
	}

	items, err := h.categoryItems(categories)
	if err != nil {
		serverError(c, err)
		return
	}
	ok(c, items)
}

func (h *KnowledgeHandler) GetCategory(c *gin.Context) {
	categoryID, okID := pathID(c, "categoryId")
	if !okID {
		return
	}

	var category model.KnowledgeCategory
	if err := h.db.Where("id = ?", categoryID).First(&category).Error; err != nil {
		serverError(c, err)
		return
	}

	var children []model.KnowledgeCategory
	if err := h.db.Where("parent_id = ?", categoryID).
		Order("sort_order ASC, id ASC").
		Find(&children).Error; err != nil {
		serverError(c, err)
		return
	}
	childItems, err := h.categoryItems(children)
	if err != nil {
		serverError(c, err)
		return
	}

	articles, err := h.findArticlesByCategory(categoryID, 20)
	if err != nil {
		serverError(c, err)
		return
	}

	breadcrumb, err := h.categoryBreadcrumb(category)
	if err != nil {
		serverError(c, err)
		return
	}

	ok(c, knowledgeCategoryDetail{
		KnowledgeCategory: category,
		Breadcrumb:        breadcrumb,
		Children:          childItems,
		Articles:          articles,
	})
}

func (h *KnowledgeHandler) ListArticles(c *gin.Context) {
	page, pageSize := pagination(c)
	query := h.db.Model(&model.KnowledgeArticle{}).Where("knowledge_articles.status = ?", statusVisible)

	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("knowledge_articles.title LIKE ? OR knowledge_articles.summary LIKE ? OR knowledge_articles.content LIKE ? OR knowledge_articles.tags LIKE ?", pattern, pattern, pattern, pattern)
	}
	if categoryID, ok := queryInt64(c, "categoryId"); ok {
		query = query.Joins("JOIN knowledge_article_categories ON knowledge_article_categories.article_id = knowledge_articles.id").
			Where("knowledge_article_categories.category_id = ?", categoryID)
	}

	var total int64
	if err := query.Distinct("knowledge_articles.id").Count(&total).Error; err != nil {
		serverError(c, err)
		return
	}

	var articles []model.KnowledgeArticle
	if err := query.Select("knowledge_articles.*").
		Group("knowledge_articles.id").
		Order("knowledge_articles.sort_order ASC, knowledge_articles.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{
		"items":    articleItems(articles, false),
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (h *KnowledgeHandler) Search(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("q"))
	if keyword == "" {
		ok(c, gin.H{"categories": []knowledgeCategoryItem{}, "articles": []knowledgeArticleItem{}})
		return
	}
	pattern := "%" + keyword + "%"

	var categories []model.KnowledgeCategory
	if err := h.db.Where("name LIKE ? OR description LIKE ?", pattern, pattern).
		Order("sort_order ASC, id ASC").
		Limit(20).
		Find(&categories).Error; err != nil {
		serverError(c, err)
		return
	}
	categoryItems, err := h.categoryItems(categories)
	if err != nil {
		serverError(c, err)
		return
	}

	var articles []model.KnowledgeArticle
	if err := h.db.Where("status = ? AND (title LIKE ? OR summary LIKE ? OR content LIKE ? OR tags LIKE ?)", statusVisible, pattern, pattern, pattern, pattern).
		Order("sort_order ASC, created_at DESC").
		Limit(20).
		Find(&articles).Error; err != nil {
		serverError(c, err)
		return
	}

	ok(c, gin.H{"categories": categoryItems, "articles": articleItems(articles, false)})
}

func (h *KnowledgeHandler) GetArticle(c *gin.Context) {
	articleID, okID := pathID(c, "articleId")
	if !okID {
		return
	}

	var article model.KnowledgeArticle
	if err := h.db.Where("id = ? AND status = ?", articleID, statusVisible).First(&article).Error; err != nil {
		serverError(c, err)
		return
	}

	var categories []model.KnowledgeCategory
	if err := h.db.Joins("JOIN knowledge_article_categories ON knowledge_article_categories.category_id = knowledge_categories.id").
		Where("knowledge_article_categories.article_id = ?", articleID).
		Order("knowledge_categories.sort_order ASC, knowledge_categories.id ASC").
		Find(&categories).Error; err != nil {
		serverError(c, err)
		return
	}

	var brand *model.KnowledgeBrand
	if article.BrandID != nil {
		var item model.KnowledgeBrand
		if err := h.db.Where("id = ?", *article.BrandID).First(&item).Error; err == nil {
			brand = &item
		} else if err != nil && err != gorm.ErrRecordNotFound {
			serverError(c, err)
			return
		}
	}

	_ = h.db.Model(&model.KnowledgeArticle{}).Where("id = ?", articleID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error

	detail := knowledgeArticleDetail{
		knowledgeArticleItem: articleItem(article, true),
		Categories:           categories,
		Brand:                brand,
	}
	detail.ViewCount++
	ok(c, detail)
}

func (h *KnowledgeHandler) categoryItems(categories []model.KnowledgeCategory) ([]knowledgeCategoryItem, error) {
	items := make([]knowledgeCategoryItem, 0, len(categories))
	for _, category := range categories {
		var childCount int64
		if err := h.db.Model(&model.KnowledgeCategory{}).Where("parent_id = ?", category.ID).Count(&childCount).Error; err != nil {
			return nil, err
		}
		var articleCount int64
		if err := h.db.Model(&model.KnowledgeArticleCategory{}).Where("category_id = ?", category.ID).Count(&articleCount).Error; err != nil {
			return nil, err
		}
		items = append(items, knowledgeCategoryItem{
			KnowledgeCategory: category,
			ChildCount:        childCount,
			ArticleCount:      articleCount,
		})
	}
	return items, nil
}

func (h *KnowledgeHandler) findArticlesByCategory(categoryID int64, limit int) ([]knowledgeArticleItem, error) {
	var articles []model.KnowledgeArticle
	query := h.db.Joins("JOIN knowledge_article_categories ON knowledge_article_categories.article_id = knowledge_articles.id").
		Where("knowledge_articles.status = ? AND knowledge_article_categories.category_id = ?", statusVisible, categoryID).
		Order("knowledge_articles.sort_order ASC, knowledge_articles.created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articleItems(articles, false), nil
}

func (h *KnowledgeHandler) categoryBreadcrumb(category model.KnowledgeCategory) ([]model.KnowledgeCategory, error) {
	result := []model.KnowledgeCategory{category}
	current := category

	for depth := 0; depth < 24 && current.ParentID != nil; depth++ {
		var parent model.KnowledgeCategory
		if err := h.db.Where("id = ?", *current.ParentID).First(&parent).Error; err != nil {
			return nil, err
		}
		result = append(result, parent)
		current = parent
	}

	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result, nil
}

func articleItems(articles []model.KnowledgeArticle, includeContent bool) []knowledgeArticleItem {
	items := make([]knowledgeArticleItem, 0, len(articles))
	for _, article := range articles {
		items = append(items, articleItem(article, includeContent))
	}
	return items
}

func articleItem(article model.KnowledgeArticle, includeContent bool) knowledgeArticleItem {
	content := ""
	if includeContent {
		content = article.Content
	}
	return knowledgeArticleItem{
		ID:          article.ID,
		Title:       article.Title,
		Summary:     article.Summary,
		Content:     content,
		Tags:        jsonStringArray(article.Tags),
		Images:      jsonStringArray(article.Images),
		VideoURL:    article.VideoURL,
		BrandID:     article.BrandID,
		PriceRange:  article.PriceRange,
		ViewCount:   article.ViewCount,
		LikeCount:   article.LikeCount,
		SortOrder:   article.SortOrder,
		PublishedAt: article.PublishedAt,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
}

func jsonStringArray(raw string) []string {
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err == nil && values != nil {
		return values
	}
	return []string{}
}

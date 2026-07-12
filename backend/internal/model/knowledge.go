package model

import "time"

type KnowledgeCategory struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id"`
	ParentID    *int64 `gorm:"column:parent_id" json:"parentId,omitempty"`
	Name        string `gorm:"column:name" json:"name"`
	Slug        string `gorm:"column:slug" json:"slug"`
	Description string `gorm:"column:description" json:"description"`
	Icon        string `gorm:"column:icon" json:"icon"`
	SortOrder   int    `gorm:"column:sort_order" json:"sortOrder"`
	HasContent  int8   `gorm:"column:has_content" json:"hasContent"`
	Timestamp
}

func (KnowledgeCategory) TableName() string {
	return "knowledge_categories"
}

type KnowledgeArticle struct {
	ID          int64      `gorm:"column:id;primaryKey" json:"id"`
	Title       string     `gorm:"column:title" json:"title"`
	Summary     string     `gorm:"column:summary" json:"summary"`
	Content     string     `gorm:"column:content" json:"content"`
	Tags        string     `gorm:"column:tags" json:"tags"`
	Images      string     `gorm:"column:images" json:"images"`
	VideoURL    string     `gorm:"column:video_url" json:"videoUrl"`
	BrandID     *int64     `gorm:"column:brand_id" json:"brandId,omitempty"`
	PriceRange  string     `gorm:"column:price_range" json:"priceRange"`
	Status      int8       `gorm:"column:status" json:"status"`
	ViewCount   int        `gorm:"column:view_count" json:"viewCount"`
	LikeCount   int        `gorm:"column:like_count" json:"likeCount"`
	SortOrder   int        `gorm:"column:sort_order" json:"sortOrder"`
	PublishedAt *time.Time `gorm:"column:published_at" json:"publishedAt,omitempty"`
	Timestamp
}

func (KnowledgeArticle) TableName() string {
	return "knowledge_articles"
}

type KnowledgeArticleCategory struct {
	ArticleID  int64     `gorm:"column:article_id;primaryKey" json:"articleId"`
	CategoryID int64     `gorm:"column:category_id;primaryKey" json:"categoryId"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (KnowledgeArticleCategory) TableName() string {
	return "knowledge_article_categories"
}

type KnowledgeBrand struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Country     string `gorm:"column:country" json:"country"`
	Story       string `gorm:"column:story" json:"story"`
	Products    string `gorm:"column:products" json:"products"`
	PriceLevel  string `gorm:"column:price_level" json:"priceLevel"`
	OfficialURL string `gorm:"column:official_url" json:"officialUrl"`
	SortOrder   int    `gorm:"column:sort_order" json:"sortOrder"`
	Timestamp
}

func (KnowledgeBrand) TableName() string {
	return "knowledge_brands"
}

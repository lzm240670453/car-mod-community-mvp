package model

import "time"

type Post struct {
	ID             int64  `gorm:"column:id;primaryKey" json:"id"`
	UserID         int64  `gorm:"column:user_id" json:"userId"`
	GarageID       *int64 `gorm:"column:garage_id" json:"garageId,omitempty"`
	Type           int8   `gorm:"column:type" json:"type"`
	Title          string `gorm:"column:title" json:"title"`
	Content        string `gorm:"column:content" json:"content"`
	VehicleModelID *int64 `gorm:"column:vehicle_model_id" json:"vehicleModelId,omitempty"`
	Status         int8   `gorm:"column:status" json:"status"`
	LikeCount      int    `gorm:"column:like_count" json:"likeCount"`
	CommentCount   int    `gorm:"column:comment_count" json:"commentCount"`
	FavoriteCount  int    `gorm:"column:favorite_count" json:"favoriteCount"`
	Timestamp
}

func (Post) TableName() string {
	return "posts"
}

type PostImage struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	PostID    int64     `gorm:"column:post_id" json:"postId"`
	ImageURL  string    `gorm:"column:image_url" json:"imageUrl"`
	SortOrder int       `gorm:"column:sort_order" json:"sortOrder"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (PostImage) TableName() string {
	return "post_images"
}

type Comment struct {
	ID       int64  `gorm:"column:id;primaryKey" json:"id"`
	PostID   int64  `gorm:"column:post_id" json:"postId"`
	UserID   int64  `gorm:"column:user_id" json:"userId"`
	ParentID *int64 `gorm:"column:parent_id" json:"parentId,omitempty"`
	Content  string `gorm:"column:content" json:"content"`
	Status   int8   `gorm:"column:status" json:"status"`
	Timestamp
}

func (Comment) TableName() string {
	return "comments"
}

type PostLike struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	PostID    int64     `gorm:"column:post_id" json:"postId"`
	UserID    int64     `gorm:"column:user_id" json:"userId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (PostLike) TableName() string {
	return "post_likes"
}

type PostFavorite struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	PostID    int64     `gorm:"column:post_id" json:"postId"`
	UserID    int64     `gorm:"column:user_id" json:"userId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (PostFavorite) TableName() string {
	return "post_favorites"
}

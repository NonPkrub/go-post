package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Content   *string    `json:"content" db:"content"`
	Published *bool      `json:"published" db:"published"`
	ViewCount int64      `json:"view_count" db:"view_count"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func NewPostModel(ID uuid.UUID, title string, content string, published bool, viewCount int64, createdAt time.Time, updatedAt time.Time, deletedAt time.Time) *Post {
	return &Post{
		ID:        ID,
		Title:     title,
		Content:   &content,
		Published: &published,
		ViewCount: viewCount,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: &deletedAt,
	}
}

type PostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostUpdateReq struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	ViewCount int       `json:"view_count"`
}

type PostRes struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
}

type PostAllReq struct {
	Published *bool     `json:"published"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"limit"`
	Count     int `json:"count"`
	TotalPage int `json:"total_page"`
}

type PostResponse struct {
	Posts     []PostRes `json:"posts"`
	Count     int       `json:"count"`
	Limit     int       `json:"limit"`
	Page      int       `json:"page"`
	TotalPage int       `json:"total_page"`
}

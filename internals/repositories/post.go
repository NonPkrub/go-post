package repositories

import (
	"fmt"
	"go-test/internals/core/domain"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) (*domain.Post, error) {
	query := "insert into post ( title, content) values ( $1, $2) RETURNING id, published,created_at"
	err := r.db.QueryRow(query, post.Title, post.Content).Scan(&post.ID, &post.Published, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepository) FindAllField(query *domain.PostAllReq, pagination *domain.Pagination) ([]*domain.Post, int64, int64, error) {
	posts := []*domain.Post{}
	var totalCount int64

	querySQL := "SELECT * FROM posts WHERE 1=1 AND deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	if query != nil {
		if query.Title != "" {
			querySQL += fmt.Sprintf(" AND title = $%d", argIndex)
			args = append(args, query.Title)
			argIndex++
		}

		if !query.CreatedAt.IsZero() {
			querySQL += fmt.Sprintf(" AND created_at >= timezone('UTC', $%d)", argIndex)
			args = append(args, query.CreatedAt)
			argIndex++
		}

		if query.Published != nil {
			querySQL += fmt.Sprintf(" AND published = $%d", argIndex)
			args = append(args, query.Published)
			argIndex++
		}
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", querySQL)
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	querySQL += fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.PageSize, offset)

	rows, err := r.db.Query(querySQL, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		post := &domain.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Published, &post.ViewCount, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
		if err != nil {
			return nil, 0, 0, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalPages := (totalCount + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return posts, totalCount, totalPages, nil
}

func (r *postRepository) FindOne(post *domain.Post) (*domain.Post, error) {
	querySQL := "SELECT id, title, content, published, view_count, created_at, updated_at, deleted_at FROM posts WHERE id = $1 AND deleted_at IS NULL"
	var postSQL domain.Post
	err := r.db.Get(&postSQL, querySQL, post.ID)
	if err != nil {
		return nil, err
	}

	if *postSQL.Published == true {
		postSQL.ViewCount += 1
		querySQL = "UPDATE post SET view_count=$1 WHERE id=$2 AND published = true AND deleted_at IS NULL"
		_, err := r.db.Exec(querySQL, postSQL.ViewCount, postSQL.ID)
		if err != nil {
			return nil, err
		}
	}

	return &postSQL, nil
}

func (r *postRepository) UpdateByID(post *domain.Post) (*domain.Post, error) {
	var fields []string
	var args []interface{}
	argIndex := 1

	if post.Title != "" {
		fields = append(fields, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, post.Title)
	}

	if post.Content != nil && *post.Content != "" {
		fields = append(fields, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, *post.Content)
	}

	if post.Published != nil {
		fields = append(fields, fmt.Sprintf("published = $%d", argIndex))
		args = append(args, *post.Published)
	}

	if post.ViewCount != 0 {
		fields = append(fields, fmt.Sprintf("view_count = $%d", argIndex))
		args = append(args, post.ViewCount)
	}

	args = append(args, post.ID)

	query := fmt.Sprintf("UPDATE posts SET %s WHERE id = $%d AND deleted_at IS NULL", strings.Join(fields, ", "), len(args))

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	updatedPost, err := r.FindOne(post)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (r *postRepository) DeleteByID(post *domain.Post) error {
	query := "UPDATE posts SET deleted_at = $1 WHERE id = $2"
	now := time.Now()
	_, err := r.db.Exec(query, now, post.ID)
	if err != nil {
		return err
	}

	return nil
}

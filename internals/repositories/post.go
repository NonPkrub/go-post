package repositories

import (
	"fmt"
	"go-test/internals/core/domain"
	"strings"

	"github.com/jmoiron/sqlx"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) (*domain.Post, error) {
	query := "insert into posts ( title, content,) values ( $1, $2) RETURNING id"
	// _, err := r.db.Exec(query, post.Title, post.Content)
	// if err != nil {
	// 	return nil, err
	// }

	// return post, nil
	err := r.db.QueryRow(query, post.Title, post.Content).Scan(&post.ID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepository) FindAllField(query *domain.PostAllReq, pagination *domain.Pagination) ([]*domain.Post, int64, int64, error) {
	posts := []*domain.Post{}
	var totalCount int64

	querySQL := "SELECT * FROM posts WHERE 1=1"
	args := []interface{}{}

	if query != nil {
		if query.Title != "" {
			querySQL += " AND title = $1"
			args = append(args, query.Title)
		}

		if !query.CreatedAt.IsZero() {
			querySQL += " AND created_at = $2"
			args = append(args, query.CreatedAt)
		}

		if query.Published != nil {
			querySQL += " AND published = $3"
			args = append(args, query.Published)
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
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Published)
		if err != nil {
			return nil, 0, 0, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalPages := (totalCount + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return nil, totalCount, totalPages, nil
}

func (r *postRepository) FindOne(post *domain.Post) (*domain.Post, error) {
	querySQL := "SELECT * FROM posts WHERE where id = $1"
	err := r.db.Get(post, querySQL, post.ID)
	if err != nil {
		return nil, err
	}
	if *post.Published == true {
		*post.ViewCount += 1
		querySQL = "UPDATE posts SET view_count=$1 WHERE id=$2 AND published = true"
		_, err := r.db.Exec(querySQL, post.ViewCount, post.ID)
		if err != nil {
			return nil, err
		}

	}

	return post, nil
}

func (r *postRepository) UpdateByID(post *domain.Post) (*domain.Post, error) {
	var fields []string
	var args []interface{}

	if post.Title != "" {
		fields = append(fields, "title = $1")
		args = append(args, post.Title)
	}

	if post.Content != nil {
		fields = append(fields, "content = $2")
		args = append(args, post.Content)
	}

	if post.Published != nil {
		fields = append(fields, "published = $3")
		args = append(args, post.Published)
	}

	if post.ViewCount != nil {
		fields = append(fields, "view_count = $4")
		args = append(args, post.ViewCount)
	}

	query := fmt.Sprintf("UPDATE posts SET %s WHERE id = $5", strings.Join(fields, ", "))
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
	query := "DELETE FROM posts WHERE id = $1"
	_, err := r.db.Exec(query, post.ID)
	if err != nil {
		return err
	}

	return nil
}

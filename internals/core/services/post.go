package services

import (
	"errors"
	"go-test/internals/core/domain"
	"go-test/internals/core/ports"
	"time"

	"github.com/google/uuid"
)

type postService struct {
	postRepository ports.PostRepository
}

func NewPostService(postRepository ports.PostRepository) *postService {
	return &postService{
		postRepository: postRepository,
	}
}

func (s *postService) Create(post *domain.PostReq) (*domain.PostRes, error) {
	if post.Title == "" {
		return nil, errors.New("title is required")
	}

	posts := &domain.Post{
		Title:   post.Title,
		Content: &post.Content,
	}

	result, err := s.postRepository.Create(posts)
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2006-01-02T15:04:05", result.CreatedAt.Format("2006-01-02T15:04:05"))

	res := &domain.PostRes{
		ID:        result.ID,
		Title:     result.Title,
		Content:   *result.Content,
		Published: *result.Published,
		CreatedAt: createdAt,
	}

	return res, nil
}

func (s *postService) UpdateByID(post *domain.PostUpdateReq) (*domain.PostRes, error) {
	posts := &domain.Post{
		ID:        post.ID,
		Title:     post.Title,
		Content:   &post.Content,
		Published: &post.Published,
	}

	res, err := s.postRepository.UpdateByID(posts)
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2006-01-02T15:04:05", res.CreatedAt.Format("2006-01-02T15:04:05"))

	return &domain.PostRes{
		ID:        res.ID,
		Title:     res.Title,
		Content:   *res.Content,
		Published: *res.Published,
		CreatedAt: createdAt,
	}, nil
}

func (s *postService) GetAll(query *domain.PostAllReq, pagination *domain.Pagination) (*domain.PostResponse, error) {
	res, count, totalPage, err := s.postRepository.FindAllField(query, pagination)
	if err != nil {
		return nil, err
	}

	response := &domain.PostResponse{
		Posts:     make([]domain.PostRes, len(res)),
		Count:     int(count),
		Limit:     pagination.PageSize,
		Page:      pagination.Page,
		TotalPage: int(totalPage),
	}

	for i, post := range res {
		createdAt, _ := time.Parse("2006-01-02T15:04:05", post.CreatedAt.Format("2006-01-02T15:04:05"))

		response.Posts[i] = domain.PostRes{
			ID:        post.ID,
			Title:     post.Title,
			Content:   *post.Content,
			Published: *post.Published,
			CreatedAt: createdAt,
		}
	}

	return response, nil

}

func (s *postService) GetByID(id string) (*domain.PostRes, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	post := &domain.Post{}
	post.ID = uuid

	res, err := s.postRepository.FindOne(post)
	if err != nil {
		return nil, err
	}

	published := *res.Published
	createdAt, _ := time.Parse("2006-01-02T15:04:05", res.CreatedAt.Format("2006-01-02T15:04:05"))

	return &domain.PostRes{
		ID:        res.ID,
		Title:     res.Title,
		Content:   *res.Content,
		Published: published,
		CreatedAt: createdAt,
	}, nil

}

func (s *postService) DeleteByID(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	post := &domain.Post{}
	post.ID = uuid
	err = s.postRepository.DeleteByID(post)
	if err != nil {
		return err
	}
	return nil
}

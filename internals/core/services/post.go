package services

import (
	"errors"
	"go-test/internals/core/domain"
	"go-test/internals/core/ports"

	"github.com/google/uuid"
)

type PostService struct {
	postRepository ports.PostRepository
}

func NewPostService(postRepository ports.PostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (s *PostService) Create(post *domain.PostReq) (*domain.PostRes, error) {
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

	res := &domain.PostRes{
		ID:        result.ID,
		Title:     result.Title,
		Content:   *result.Content,
		Published: *result.Published,
		CreatedAt: result.CreatedAt,
	}

	return res, nil
}

func (s *PostService) UpdateByID(post *domain.PostUpdateReq) (*domain.PostRes, error) {
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

	return &domain.PostRes{
		ID:        res.ID,
		Title:     res.Title,
		Content:   *res.Content,
		Published: *res.Published,
	}, nil
}

func (s *PostService) GetAll(query *domain.PostAllReq, pagination *domain.Pagination) (*domain.PostResponse, error) {
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
		response.Posts[i] = domain.PostRes{
			ID:        post.ID,
			Title:     post.Title,
			Content:   *post.Content,
			Published: *post.Published,
		}
	}

	return response, nil

}

func (s *PostService) GetByID(id string) (*domain.PostRes, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	postID := &domain.Post{
		ID: uuid,
	}

	res, err := s.postRepository.FindOne(postID)
	if err != nil {
		return nil, err
	}

	return &domain.PostRes{
		ID:        res.ID,
		Title:     res.Title,
		Content:   *res.Content,
		Published: *res.Published,
	}, nil

}

func (s *PostService) DeleteByID(id string) error {
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

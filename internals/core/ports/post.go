package ports

import (
	"go-test/internals/core/domain"

	"github.com/gofiber/fiber/v2"
)

type PostService interface {
	Create(post *domain.PostReq) (*domain.PostRes, error)
	GetAll(query *domain.PostAllReq, pagination *domain.Pagination) (*domain.PostResponse, error)
	GetByID(id string) (*domain.PostRes, error)
	UpdateByID(post *domain.PostUpdateReq) (*domain.PostRes, error)
	DeleteByID(id string) error
}

type PostRepository interface {
	Create(post *domain.Post) (*domain.Post, error)
	FindAllField(query *domain.PostAllReq, pagination *domain.Pagination) ([]*domain.Post, int64, int64, error)
	FindOne(post *domain.Post) (*domain.Post, error)
	UpdateByID(post *domain.Post) (*domain.Post, error)
	DeleteByID(post *domain.Post) error
}

type PostHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	UpdateByID(c *fiber.Ctx) error
	DeleteByID(c *fiber.Ctx) error
}

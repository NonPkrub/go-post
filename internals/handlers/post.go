package handlers

import (
	"go-test/internals/core/domain"
	"go-test/internals/core/ports"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PostHandler struct {
	postService ports.PostService
}

func NewPostHandler(postService ports.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) Create(c *fiber.Ctx) error {
	var form domain.PostReq
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	res, err := h.postService.Create(&form)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if res.Title == "" {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"message": "title is required",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"data": res,
	})

}

func (h *PostHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	pagination := domain.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	publishedStr := c.Query("published", "")
	var published *bool
	if publishedStr == "" {
		published = nil
	} else {
		var err error
		publish, err := strconv.ParseBool(publishedStr)
		published = &publish
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid value for 'published'",
			})
		}

	}

	title := c.Query("title", "")

	location, _ := time.LoadLocation("Asia/Bangkok")
	createdAtStr := c.Query("created_at", "")
	var createdAt time.Time
	if createdAtStr != "" {
		var err error
		createdAt, err = time.ParseInLocation(time.RFC3339, createdAtStr, location)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid value for 'created_at'.",
			})
		}

	}
	reqParams := domain.PostAllReq{
		Published: published,
		Title:     title,
		CreatedAt: createdAt,
	}

	res, err := h.postService.GetAll(&reqParams, &pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.postService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *PostHandler) UpdateByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var form domain.PostUpdateReq
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	form.ID = uuid

	res, err := h.postService.UpdateByID(&form)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": res,
	})
}

func (h *PostHandler) DeleteByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.postService.DeleteByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "deleted successfully",
	})
}

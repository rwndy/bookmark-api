package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/rwndy/bookmark-api/internal/domain"
	"github.com/rwndy/bookmark-api/internal/handler/dto"
	"github.com/rwndy/bookmark-api/pkg/response"
)

type BookmarkHandler struct {
	service domain.BookmarkService
}

func NewBookmarkHandler(service domain.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{service: service}
}

func (h *BookmarkHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.CreateBookmarkRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, 400, "invalid request body")
	}

	if err := ValidateStruct(&req); err != nil {
		return response.Fail(c, 400, err.Error())
	}

	bookmark, err := h.service.CreateBookmark(userID, req.URL, req.Title, req.Tags)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return response.Fail(c, appErr.Code, appErr.Message)
		}
		return response.Fail(c, 500, "internal server error")
	}

	return response.Created(c, "bookmark created", bookmark)
}

func (h *BookmarkHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	bookmarks, err := h.service.GetUserBookmarks(userID)
	if err != nil {
		return response.Fail(c, 500, "failed to fetch bookmarks")
	}

	return response.OK(c, "bookmarks fetched", bookmarks)
}

func (h *BookmarkHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Fail(c, 400, "invalid bookmark id")
	}

	var req dto.UpdateBookmarkRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, 400, "invalid request body")
	}

	if err := ValidateStruct(&req); err != nil {
		return response.Fail(c, 400, err.Error())
	}

	bookmark, err := h.service.UpdateBookmark(uint(id), userID, req.URL, req.Title, req.Tags)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return response.Fail(c, appErr.Code, appErr.Message)
		}
		return response.Fail(c, 500, "internal server error")
	}

	return response.OK(c, "bookmark updated", bookmark)
}

func (h *BookmarkHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Fail(c, 400, "invalid bookmark id")
	}

	if err := h.service.DeleteBookmark(uint(id), userID); err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return response.Fail(c, appErr.Code, appErr.Message)
		}
		return response.Fail(c, 500, "internal server error")
	}

	return c.SendStatus(204)
}
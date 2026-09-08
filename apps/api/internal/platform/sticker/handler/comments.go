package handler

import (
	"strconv"
	"strings"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/pkg/errors"
	"kun-galgame-sticker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) ListComments(c fiber.Ctx) error {
	packID, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	page, appErr := h.svc.PackComments(c.Context(), packID, strings.TrimSpace(c.Query("after")), viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, page)
}

func (h *Handler) AddComment(c fiber.Ctx) error {
	packID, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	req, appErr := parseBody[dto.CommentRequest](c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	comment, appErr := h.svc.AddComment(c.Context(), packID, viewer(c), req.Body, req.ReplyTo)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, comment)
}

func (h *Handler) PatchComment(c fiber.Ctx) error {
	postID, appErr := commentIDParam(c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	req, appErr := parseBody[dto.CommentRequest](c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	comment, appErr := h.svc.EditComment(c.Context(), postID, viewer(c), req.Body)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, comment)
}

func (h *Handler) DeleteComment(c fiber.Ctx) error {
	postID, appErr := commentIDParam(c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	if appErr := h.svc.DeleteComment(c.Context(), postID, viewer(c)); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, fiber.Map{"deleted": true})
}

// Comment ids are community's, not this site's: plain positive integers.
func commentIDParam(c fiber.Ctx) (int64, *errors.AppError) {
	id, err := strconv.ParseInt(strings.TrimSpace(c.Params("commentId")), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.ErrInvalidParams("commentId must be a positive integer")
	}
	return id, nil
}

package handler

import (
	"fmt"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/pkg/errors"
	"kun-galgame-sticker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) GetSticker(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "stickerId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	sticker, appErr := h.svc.GetSticker(id, viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, sticker)
}

func (h *Handler) DownloadSticker(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "stickerId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	sticker, appErr := h.svc.StickerForDownload(id, viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	body, size, appErr := h.svc.OpenImage(c.Context(), sticker.ImageHash)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	// No defer Close here: fasthttp reads the stream after this handler
	// returns and closes it itself. Closing early sent an empty response.
	h.svc.BumpDownload(sticker.PackID)
	c.Set(fiber.HeaderContentType, "image/webp")
	c.Set(fiber.HeaderContentDisposition,
		fmt.Sprintf(`attachment; filename="%s.webp"`, sticker.ID))
	if size > 0 {
		return c.SendStream(body, int(size))
	}
	return c.SendStream(body)
}

func (h *Handler) UploadImage(c fiber.Ctx) error {
	packID, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	header, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, errors.ErrInvalidParams("a multipart file field is required"))
	}
	file, err := header.Open()
	if err != nil {
		return response.Error(c, errors.ErrInvalidParams("could not read the uploaded file"))
	}
	defer file.Close()

	result, appErr := h.svc.UploadImage(
		c.Context(), packID, viewer(c), header.Filename, file, header.Size,
	)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, result)
}

func (h *Handler) AddSticker(c fiber.Ctx) error {
	packID, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	req, appErr := parseBody[dto.CreateStickerRequest](c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	sticker, appErr := h.svc.AddSticker(packID, viewer(c), req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, sticker)
}

func (h *Handler) PatchSticker(c fiber.Ctx) error {
	packID, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	stickerID, appErr := uuidParam(c, "stickerId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	req, appErr := parseBody[dto.PatchStickerRequest](c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	sticker, appErr := h.svc.PatchSticker(packID, stickerID, viewer(c), req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, sticker)
}

func (h *Handler) DeleteSticker(c fiber.Ctx) error {
	packID, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	stickerID, appErr := uuidParam(c, "stickerId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	if appErr := h.svc.DeleteSticker(packID, stickerID, viewer(c)); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, nil)
}

func (h *Handler) ReorderStickers(c fiber.Ctx) error {
	packID, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	req, appErr := parseBody[dto.ReorderRequest](c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	if appErr := h.svc.Reorder(packID, viewer(c), req.StickerIDs); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, nil)
}

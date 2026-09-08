package handler

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/pkg/errors"
	"kun-galgame-sticker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) ListPacks(c fiber.Ctx) error {
	page, appErr := h.svc.List(c.Context(), listQuery(c), viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, page)
}

func (h *Handler) ListUserPacks(c fiber.Ctx) error {
	uid, err := strconv.Atoi(c.Params("uid"))
	if err != nil || uid <= 0 {
		return response.Error(c, errors.ErrInvalidParams("uid must be a positive integer"))
	}
	query := listQuery(c)
	query.OwnerUID = uid

	page, appErr := h.svc.List(c.Context(), query, viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, page)
}

func (h *Handler) ListMyPacks(c fiber.Ctx) error {
	v := viewer(c)
	query := listQuery(c)
	query.OwnerUID = v.UID
	query.Sort = dto.SortUpdated
	query.Rating = dto.RatingFilterAll

	page, appErr := h.svc.List(c.Context(), query, v)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, page)
}

func (h *Handler) GetPack(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	detail, appErr := h.svc.Get(c.Context(), id, viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, detail)
}

func (h *Handler) CreatePack(c fiber.Ctx) error {
	req, appErr := parseBody[dto.CreatePackRequest](c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	pack, appErr := h.svc.CreatePack(c.Context(), viewer(c), req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, pack)
}

func (h *Handler) PatchPack(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	req, appErr := parseBody[dto.PatchPackRequest](c)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	pack, appErr := h.svc.PatchPack(c.Context(), id, viewer(c), req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, pack)
}

func (h *Handler) DeletePack(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	if appErr := h.svc.DeletePack(id, viewer(c)); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, nil)
}

func (h *Handler) Publish(c fiber.Ctx) error   { return h.setPublished(c, true) }
func (h *Handler) Unpublish(c fiber.Ctx) error { return h.setPublished(c, false) }

func (h *Handler) setPublished(c fiber.Ctx, publish bool) error {
	id, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	pack, appErr := h.svc.SetPublished(c.Context(), id, viewer(c), publish)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, pack)
}

func (h *Handler) DownloadPack(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return response.Error(c, appErr)
	}
	pack, stickers, appErr := h.svc.PackForDownload(id, viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}

	c.Set(fiber.HeaderContentType, "application/zip")
	c.Set(fiber.HeaderContentDisposition,
		fmt.Sprintf(`attachment; filename="pack-%s.zip"`, pack.ID))
	h.svc.BumpDownload(pack.ID)

	// The stream writer runs after this handler returns, so the request context
	// is already cancelled by then; WriteZip carries its own deadline.
	return c.SendStreamWriter(func(w *bufio.Writer) {
		if err := h.svc.WriteZip(context.Background(), w, stickers); err != nil {
			slog.Error("pack zip", "pack", pack.ID, "error", err)
		}
		_ = w.Flush()
	})
}

func (h *Handler) ListTags(c fiber.Ctx) error {
	tags, appErr := h.svc.PopularTags(maxTagsShown)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, tags)
}

// AvatarPool serves the ecosystem-wide default-avatar manifest. Public and
// unauthenticated on purpose: it is a list of URLs that are already public,
// and the sites that consume it fetch it from their own servers at boot, not
// from a browser.
func (h *Handler) AvatarPool(c fiber.Ctx) error {
	pool, appErr := h.svc.AvatarPool()
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, pool)
}

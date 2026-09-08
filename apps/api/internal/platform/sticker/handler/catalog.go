package handler

import (
	"strings"

	"kun-galgame-sticker-api/pkg/errors"
	"kun-galgame-sticker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// SearchCatalogWorks and CatalogWorkRoster relay the two catalog reads the
// editor needs. They are relays, not proxies: the browser sends a query, this
// site sends its application key, and only the fields the picker renders come
// back.
func (h *Handler) SearchCatalogWorks(c fiber.Ctx) error {
	q := truncate(strings.TrimSpace(c.Query("q")), maxSearchLen)
	works, appErr := h.svc.SearchWorks(c.Context(), q)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, fiber.Map{"works": works})
}

func (h *Handler) CatalogWorkRoster(c fiber.Ctx) error {
	workID := catalogID(c.Params("workId"))
	if workID == 0 {
		return response.Error(c, errors.ErrInvalidParams("workId must be a catalog id"))
	}
	characters, appErr := h.svc.WorkRoster(c.Context(), workID)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, fiber.Map{"characters": characters})
}

func (h *Handler) GetCharacter(c fiber.Ctx) error {
	characterID := catalogID(c.Params("characterId"))
	if characterID == 0 {
		return response.Error(c, errors.ErrInvalidParams("characterId must be a catalog id"))
	}
	page, appErr := h.svc.CharacterPage(c.Context(), characterID, viewer(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, page)
}

package service

import (
	"context"
	"strings"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/internal/platform/sticker/model"
	"kun-galgame-sticker-api/internal/platform/sticker/repository"
	"kun-galgame-sticker-api/pkg/errors"
)

// searchLimit is per lane, not overall: the palette shows two short lists and
// a "see all" that leads to the full discovery page.
const searchLimit = 6

// QuickSearch backs the command palette. It answers from this site's own data
// -- packs and the characters those packs are tagged with -- rather than from
// catalog, because every hit has to lead to a page that exists here. A catalog
// character with no stickers on this site would be a dead end dressed up as a
// result.
func (s *Service) QuickSearch(ctx context.Context, query string, viewer Viewer) (*dto.SearchResults, *errors.AppError) {
	query = strings.TrimSpace(query)
	out := &dto.SearchResults{
		Packs:      []dto.Pack{},
		Characters: []dto.CatalogCharacter{},
	}
	if query == "" {
		return out, nil
	}

	rows, _, err := s.packs.List(repository.ListParams{
		Statuses: []int16{model.PackPublished},
		Search:   query,
		Order:    orderFor(dto.SortHot),
		Limit:    searchLimit,
	})
	if err != nil {
		return nil, errors.ErrInternal("failed to search packs")
	}
	packs, appErr := s.hydrate(ctx, rows)
	if appErr != nil {
		return nil, appErr
	}
	out.Packs = packs

	hits, err := s.stickers.SearchCharacters(query, searchLimit)
	if err != nil {
		return nil, errors.ErrInternal("failed to search characters")
	}
	for _, hit := range hits {
		out.Characters = append(out.Characters, dto.CatalogCharacter{
			ID:           hit.CatalogCharacterID,
			Name:         decodeML(hit.CatalogCharacterName),
			ImageURL:     hit.CatalogCharacterImage,
			WorkName:     decodeML(hit.CatalogWorkName),
			StickerCount: hit.StickerCount,
		})
	}
	return out, nil
}

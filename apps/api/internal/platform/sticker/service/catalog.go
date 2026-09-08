package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/internal/platform/sticker/model"
	"kun-galgame-sticker-api/pkg/catalogclient"
	"kun-galgame-sticker-api/pkg/errors"

	"gorm.io/datatypes"
)

// catalogName wraps the shared fold into this package's DTO type. The mapping
// itself lives in catalogclient because the backfill command needs the same one.
func catalogName(displayName string, localized map[string]catalogclient.LocalizedText, latin *string) dto.MultilingualText {
	return catalogclient.LocalizedName(displayName, localized, latin)
}

func encodeML(in dto.MultilingualText) datatypes.JSON {
	raw, err := json.Marshal(in)
	if err != nil {
		return datatypes.JSON("{}")
	}
	return datatypes.JSON(raw)
}

// imageURL is the single gate on catalog imagery. catalog rates an image
// safe | suggestive | explicit, but leaves nearly every row unassessed, so
// this can only drop what is actually marked explicit -- it is a backstop, not
// a guarantee. The reliable signal is the work's content_rating, which travels
// to the UI as a badge instead.
func imageURL(img *catalogclient.Image) string {
	if img == nil {
		return ""
	}
	if img.Sexual != nil && *img.Sexual == "explicit" {
		return ""
	}
	return img.URL
}

func parseCatalogID(raw string) *int64 {
	id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || id <= 0 {
		return nil
	}
	return &id
}

// workDTO and characterDTO render the snapshot the database holds. They never
// reach for catalog: a list page must stay renderable while it is down.
// workSnapshot is what a row stores about a catalog work. It travels as one
// value because the caller writes all of it or none of it.
type workSnapshot struct {
	ID     *int64
	Name   datatypes.JSON
	Cover  string
	Rating string
}

func emptyWorkSnapshot() workSnapshot {
	return workSnapshot{Name: datatypes.JSON("{}")}
}

func workDTO(id *int64, name datatypes.JSON, cover, rating string) *dto.CatalogWork {
	if id == nil {
		return nil
	}
	return &dto.CatalogWork{
		ID:            *id,
		Name:          decodeML(name),
		CoverURL:      cover,
		ContentRating: rating,
	}
}

func characterDTO(id *int64, name datatypes.JSON, image string) *dto.CatalogCharacter {
	if id == nil {
		return nil
	}
	return &dto.CatalogCharacter{
		ID:       *id,
		Name:     decodeML(name),
		ImageURL: image,
	}
}

// resolveWork turns a catalog work id into the snapshot stored on a row. A nil
// or zero id clears the link; an id catalog does not know is a client error,
// not a silently dropped field.
func (s *Service) resolveWork(ctx context.Context, id *int64) (workSnapshot, *errors.AppError) {
	if id == nil || *id <= 0 {
		return emptyWorkSnapshot(), nil
	}
	if !s.catalog.Configured() {
		return workSnapshot{}, errors.ErrCatalogUnavailable()
	}
	work, err := s.catalog.Work(ctx, *id)
	if err != nil {
		if catalogclient.Missing(err) {
			return workSnapshot{}, errors.ErrInvalidParams("no such catalog work")
		}
		return workSnapshot{}, errors.ErrCatalogUnavailable()
	}
	name := catalogName(work.DisplayName, work.Localized, work.Latin)
	return workSnapshot{
		ID:     id,
		Name:   encodeML(name),
		Cover:  imageURL(work.Cover),
		Rating: work.ContentRating,
	}, nil
}

func (s *Service) resolveCharacter(ctx context.Context, id *int64) (*int64, datatypes.JSON, string, *errors.AppError) {
	if id == nil || *id <= 0 {
		return nil, datatypes.JSON("{}"), "", nil
	}
	if !s.catalog.Configured() {
		return nil, nil, "", errors.ErrCatalogUnavailable()
	}
	character, err := s.catalog.CharacterDetail(ctx, *id)
	if err != nil {
		if catalogclient.Missing(err) {
			return nil, nil, "", errors.ErrInvalidParams("no such catalog character")
		}
		return nil, nil, "", errors.ErrCatalogUnavailable()
	}
	name := catalogName(character.DisplayName, character.Localized, character.Latin)
	return id, encodeML(name), imageURL(character.Image), nil
}

// SearchWorks and WorkRoster back the editor's two pickers. They are thin
// passthroughs: the browser must not hold the application key, so every
// catalog read for a logged-in author goes through this site's own API.
func (s *Service) SearchWorks(ctx context.Context, q string) ([]dto.CatalogWork, *errors.AppError) {
	q = strings.TrimSpace(q)
	if q == "" {
		return []dto.CatalogWork{}, nil
	}
	if !s.catalog.Configured() {
		return nil, errors.ErrCatalogUnavailable()
	}
	works, err := s.catalog.SearchWorks(ctx, q, 0)
	if err != nil {
		return nil, errors.ErrCatalogUnavailable()
	}
	out := make([]dto.CatalogWork, 0, len(works))
	for _, work := range works {
		id := parseCatalogID(work.ID)
		if id == nil {
			continue
		}
		out = append(out, dto.CatalogWork{
			ID:            *id,
			Name:          catalogName(work.DisplayName, work.Localized, work.Latin),
			CoverURL:      imageURL(work.Cover),
			ReleaseDate:   work.ReleaseDate,
			Medium:        work.Medium,
			ContentRating: work.ContentRating,
		})
	}
	return out, nil
}

func (s *Service) WorkRoster(ctx context.Context, workID int64) ([]dto.CatalogCharacter, *errors.AppError) {
	if !s.catalog.Configured() {
		return nil, errors.ErrCatalogUnavailable()
	}
	roster, err := s.catalog.WorkCharacters(ctx, workID)
	if err != nil {
		if catalogclient.Missing(err) {
			return nil, errors.ErrInvalidParams("no such catalog work")
		}
		return nil, errors.ErrCatalogUnavailable()
	}
	out := make([]dto.CatalogCharacter, 0, len(roster))
	for _, character := range roster {
		id := parseCatalogID(character.ID)
		if id == nil {
			continue
		}
		out = append(out, dto.CatalogCharacter{
			ID:         *id,
			Name:       catalogName(character.DisplayName, character.Localized, character.Latin),
			ImageURL:   imageURL(character.Image),
			RosterRole: character.RosterRole,
		})
	}
	return out, nil
}

// CharacterPage is the public character route: catalog's profile on the left,
// this site's stickers of that character on the right. The catalog half is
// optional -- when the upstream is down the page still lists the stickers,
// using the name snapshot they carry.
func (s *Service) CharacterPage(ctx context.Context, characterID int64, viewer Viewer) (*dto.CharacterPage, *errors.AppError) {
	rows, err := s.stickers.ByCatalogCharacter(characterID, characterStickerLimit)
	if err != nil {
		return nil, errors.ErrInternal("failed to load stickers")
	}
	packs, appErr := s.packsForStickers(ctx, rows, viewer)
	if appErr != nil {
		return nil, appErr
	}

	out := &dto.CharacterPage{
		Character: dto.CatalogCharacter{ID: characterID},
		Stickers:  make([]dto.CharacterSticker, 0, len(rows)),
		Packs:     packs,
	}
	for _, row := range rows {
		pack, ok := packs[row.PackID.String()]
		if !ok {
			continue
		}
		main, thumb := s.urls(row.ImageHash)
		out.Stickers = append(out.Stickers, dto.CharacterSticker{
			ID:       row.ID.String(),
			PackID:   pack.ID,
			ImageURL: main,
			ThumbURL: thumb,
			Width:    row.Width,
			Height:   row.Height,
		})
		if len(out.Character.Name) == 0 {
			out.Character.Name = decodeML(row.CatalogCharacterName)
			out.Character.ImageURL = row.CatalogCharacterImage
		}
	}

	if s.catalog.Configured() {
		if detail, err := s.catalog.CharacterDetail(ctx, characterID); err == nil {
			out.Character.Name = catalogName(detail.DisplayName, detail.Localized, detail.Latin)
			out.Character.ImageURL = imageURL(detail.Image)
			out.Character.Gender = detail.Gender
			out.Character.Birthday = detail.Birthday
			out.Character.BloodType = detail.BloodType
			out.Character.Traits = traitDTOs(detail.Traits)
			out.Character.Aliases = aliasNames(detail.Aliases)
			out.Profile = true
		} else if catalogclient.Missing(err) && len(out.Stickers) == 0 {
			return nil, errors.ErrCharacterNotFound()
		}
		if works, err := s.catalog.CharacterAppearances(ctx, characterID, 12); err == nil {
			out.Appearances = s.appearanceDTOs(ctx, works)
		}
	}

	if len(out.Stickers) == 0 && !out.Profile {
		return nil, errors.ErrCharacterNotFound()
	}
	return out, nil
}

// packsForStickers loads the packs those stickers belong to and drops the ones
// this viewer may not see, so an unpublished pack never leaks a sticker
// through the character page.
func (s *Service) packsForStickers(ctx context.Context, rows []model.Sticker, viewer Viewer) (map[string]dto.Pack, *errors.AppError) {
	ids := make([]string, 0, len(rows))
	seen := map[string]bool{}
	for _, row := range rows {
		key := row.PackID.String()
		if !seen[key] {
			seen[key] = true
			ids = append(ids, key)
		}
	}
	packs, err := s.packs.ByIDs(ids)
	if err != nil {
		return nil, errors.ErrInternal("failed to load packs")
	}
	visible := make([]model.Pack, 0, len(packs))
	for _, pack := range packs {
		if pack.Status == model.PackPublished || viewer.canSeeUnpublished(&pack) {
			visible = append(visible, pack)
		}
	}
	hydrated, appErr := s.hydrate(ctx, visible)
	if appErr != nil {
		return nil, appErr
	}
	out := make(map[string]dto.Pack, len(hydrated))
	for _, pack := range hydrated {
		out[pack.ID] = pack
	}
	return out, nil
}

// traitDTOs keeps the localized trait and group names catalog ships alongside
// the English ones, so a Chinese reader gets 马尾辫 / 毛发 rather than
// Ponytail / Hair. Sexual traits are dropped: this is a sticker site.
func traitDTOs(in *[]catalogclient.Trait) []dto.CatalogTrait {
	if in == nil {
		return nil
	}
	out := make([]dto.CatalogTrait, 0, len(*in))
	for _, trait := range *in {
		if trait.IsSexual {
			continue
		}
		out = append(out, dto.CatalogTrait{
			Name:  catalogName(trait.DisplayName, trait.Localized, nil),
			Group: catalogName(trait.Group, trait.GroupLocalized, nil),
		})
	}
	return out
}

func aliasNames(in *[]catalogclient.Alias) []string {
	if in == nil {
		return nil
	}
	out := make([]string, 0, len(*in))
	for _, alias := range *in {
		if name := strings.TrimSpace(alias.Name); name != "" {
			out = append(out, name)
		}
	}
	return out
}

// appearanceDTOs fills in the covers the appearances lane leaves null with one
// batch lookup, so the list reads as a shelf of games rather than a row of
// empty boxes. A failed lookup costs the covers, not the list.
func (s *Service) appearanceDTOs(ctx context.Context, in []catalogclient.Appearance) []dto.CatalogWork {
	ids := make([]string, 0, len(in))
	for _, appearance := range in {
		ids = append(ids, appearance.Work.ID)
	}
	covers, err := s.catalog.WorksByIDs(ctx, ids)
	if err != nil {
		covers = nil
	}

	out := make([]dto.CatalogWork, 0, len(in))
	for _, appearance := range in {
		id := parseCatalogID(appearance.Work.ID)
		if id == nil {
			continue
		}
		work := appearance.Work
		if full, ok := covers[work.ID]; ok {
			work = full
		}
		out = append(out, dto.CatalogWork{
			ID:            *id,
			Name:          catalogName(work.DisplayName, work.Localized, work.Latin),
			CoverURL:      imageURL(work.Cover),
			ReleaseDate:   work.ReleaseDate,
			Medium:        work.Medium,
			ContentRating: work.ContentRating,
		})
	}
	return out
}

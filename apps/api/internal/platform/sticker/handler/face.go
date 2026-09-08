package handler

import (
	"strconv"
	"strings"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/internal/platform/sticker/repository"
	"kun-galgame-sticker-api/pkg/errors"
	"kun-galgame-sticker-api/pkg/problem"

	"github.com/gofiber/fiber/v3"
)

// Handlers for the public developer-platform face. Two things separate them
// from the site's own handlers: the body is bare JSON rather than the house
// {code, message, data} envelope, and failures are RFC 9457 problem documents.
// Authentication is not among them -- the gateway terminates it, and by the
// time a request arrives here it has already been keyed, budgeted and metered.

const (
	faceDefaultLimit = 20
	faceMaxLimit     = 50
	faceMaxPage      = 1000
)

// faceFail translates a service error into a problem document. The service
// speaks the site's error codes because it is shared with the site; this is the
// one place that mapping is written down.
func faceFail(c fiber.Ctx, err *errors.AppError) error {
	switch err.StatusCode {
	case fiber.StatusNotFound:
		return problem.Write(c, problem.CodeNotFound, err.Message)
	case fiber.StatusBadRequest:
		return problem.Write(c, problem.CodeInvalidParameter, err.Message)
	case fiber.StatusServiceUnavailable:
		return problem.Write(c, problem.CodeServiceUnavailable, err.Message)
	default:
		return problem.Write(c, problem.CodeInternalError, err.Message)
	}
}

// faceFault is a rejected parameter that has not been written yet. The
// helpers below used to return the result of problem.Write, which is nil on
// success -- so `if err != nil` never fired, the handler ran the query anyway
// and the second write replaced the problem document with a success body under
// a 400 status line. A pointer that is either nil or a fault cannot do that.
type faceFault struct{ code, detail string }

func (f *faceFault) write(c fiber.Ctx) error { return problem.Write(c, f.code, f.detail) }

func invalid(detail string) *faceFault {
	return &faceFault{code: problem.CodeInvalidParameter, detail: detail}
}

// faceLimit rejects rather than clamps, unlike the site's own handlers: a
// third party writing a paging loop needs to be told its page size was
// refused, where a browser just wants the biggest page it may have.
// LIMIT_TOO_LARGE is catalog's code for exactly this.
func faceLimit(c fiber.Ctx) (int, *faceFault) {
	raw := strings.TrimSpace(c.Query("limit"))
	if raw == "" {
		return faceDefaultLimit, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, invalid("limit must be a positive integer")
	}
	if n > faceMaxLimit {
		return 0, &faceFault{
			code:   problem.CodeLimitTooLarge,
			detail: "limit must be at most " + strconv.Itoa(faceMaxLimit),
		}
	}
	return n, nil
}

func facePage(c fiber.Ctx) (int, *faceFault) {
	raw := strings.TrimSpace(c.Query("page"))
	if raw == "" {
		return 1, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > faceMaxPage {
		return 0, invalid("page must be between 1 and " + strconv.Itoa(faceMaxPage))
	}
	return n, nil
}

// faceBool is a closed vocabulary, as catalog's flags are: absent means the
// default, and anything other than true or false is a client error rather than
// a silent false.
func faceBool(c fiber.Ctx, name string) (value, set bool, bad *faceFault) {
	switch strings.TrimSpace(c.Query(name)) {
	case "":
		return false, false, nil
	case "true":
		return true, true, nil
	case "false":
		return false, true, nil
	}
	return false, false, invalid(name + " must be true or false")
}

func faceCatalogID(c fiber.Ctx, name string) (int64, *faceFault) {
	id, err := strconv.ParseInt(strings.TrimSpace(c.Params(name)), 10, 64)
	if err != nil || id <= 0 {
		return 0, invalid(name + " must be a positive integer")
	}
	return id, nil
}

func faceSort(c fiber.Ctx) (string, *faceFault) {
	raw := strings.TrimSpace(c.Query("sort"))
	if raw != "" && raw != dto.SortNew && raw != dto.SortHot {
		return "", invalid("sort must be new or hot")
	}
	return sortOf(raw), nil
}

// faceRatingFilter turns the nsfw flag into this site's rating filter. Absent
// or false hides r18 packs, matching catalog's nsfw parameter exactly.
func faceRatingFilter(c fiber.Ctx) (string, *faceFault) {
	nsfw, _, bad := faceBool(c, "nsfw")
	if bad != nil {
		return "", bad
	}
	if nsfw {
		return dto.RatingFilterAll, nil
	}
	return dto.RatingFilterSFW, nil
}

func (h *Handler) FaceListPacks(c fiber.Ctx) error {
	limit, bad := faceLimit(c)
	if bad != nil {
		return bad.write(c)
	}
	page, bad := facePage(c)
	if bad != nil {
		return bad.write(c)
	}
	return h.faceListPacks(c, limit, page, 0)
}

// faceListPacks is shared by /packs and /works/{id}/packs, which differ only in
// where the work filter comes from.
func (h *Handler) faceListPacks(c fiber.Ctx, limit, page int, work int64) error {
	sort, bad := faceSort(c)
	if bad != nil {
		return bad.write(c)
	}
	rating, bad := faceRatingFilter(c)
	if bad != nil {
		return bad.write(c)
	}
	official, officialSet, bad := faceBool(c, "official")
	if bad != nil {
		return bad.write(c)
	}
	linked, linkedSet, bad := faceBool(c, "linked")
	if bad != nil {
		return bad.write(c)
	}

	query := dto.ListQuery{
		Page:         page,
		Limit:        limit,
		Sort:         sort,
		Search:       truncate(strings.TrimSpace(c.Query("q")), maxSearchLen),
		Tag:          repository.Slugify(c.Query("tag")),
		Rating:       rating,
		OfficialOnly: officialSet && official,
		LinkedOnly:   linkedSet && linked,
	}
	if work > 0 {
		query.AnyCatalogWorkID = work
	} else {
		query.AnyCatalogWorkID = catalogID(c.Query("work"))
	}
	out, appErr := h.svc.FaceListPacks(c.Context(), query)
	if appErr != nil {
		return faceFail(c, appErr)
	}
	return c.JSON(out)
}

func (h *Handler) FaceGetPack(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "packId")
	if appErr != nil {
		return problem.Write(c, problem.CodeInvalidParameter, appErr.Message)
	}
	out, svcErr := h.svc.FaceGetPack(c.Context(), id)
	if svcErr != nil {
		return faceFail(c, svcErr)
	}
	return c.JSON(out)
}

func (h *Handler) FaceGetSticker(c fiber.Ctx) error {
	id, appErr := uuidParam(c, "stickerId")
	if appErr != nil {
		return problem.Write(c, problem.CodeInvalidParameter, appErr.Message)
	}
	out, svcErr := h.svc.FaceGetSticker(id)
	if svcErr != nil {
		return faceFail(c, svcErr)
	}
	return c.JSON(out)
}

func (h *Handler) FaceListCharacters(c fiber.Ctx) error {
	params, bad := faceIndexParams(c)
	if bad != nil {
		return bad.write(c)
	}
	params.Work = catalogID(c.Query("work"))
	out, appErr := h.svc.FaceCharacters(params)
	if appErr != nil {
		return faceFail(c, appErr)
	}
	return c.JSON(out)
}

func (h *Handler) FaceGetCharacter(c fiber.Ctx) error {
	id, bad := faceCatalogID(c, "characterId")
	if bad != nil {
		return bad.write(c)
	}
	out, appErr := h.svc.FaceCharacter(id)
	if appErr != nil {
		return faceFail(c, appErr)
	}
	return c.JSON(out)
}

func (h *Handler) FaceCharacterStickers(c fiber.Ctx) error {
	id, bad := faceCatalogID(c, "characterId")
	if bad != nil {
		return bad.write(c)
	}
	params, bad := faceIndexParams(c)
	if bad != nil {
		return bad.write(c)
	}
	out, appErr := h.svc.FaceCharacterStickers(id, params.Offset, params.Limit)
	if appErr != nil {
		return faceFail(c, appErr)
	}
	return c.JSON(out)
}

func (h *Handler) FaceListWorks(c fiber.Ctx) error {
	params, bad := faceIndexParams(c)
	if bad != nil {
		return bad.write(c)
	}
	out, appErr := h.svc.FaceWorks(params)
	if appErr != nil {
		return faceFail(c, appErr)
	}
	return c.JSON(out)
}

// FaceWorkPacks is /packs?work= under the path a caller holding a catalog work
// id reaches for. "About this work" means the pack declares it or holds a
// sticker of it -- the same relation the work index counts.
func (h *Handler) FaceWorkPacks(c fiber.Ctx) error {
	id, bad := faceCatalogID(c, "workId")
	if bad != nil {
		return bad.write(c)
	}
	limit, bad := faceLimit(c)
	if bad != nil {
		return bad.write(c)
	}
	page, bad := facePage(c)
	if bad != nil {
		return bad.write(c)
	}
	return h.faceListPacks(c, limit, page, id)
}

func (h *Handler) FaceListTags(c fiber.Ctx) error {
	limit, bad := faceLimit(c)
	if bad != nil {
		return bad.write(c)
	}
	out, appErr := h.svc.FaceTags(limit)
	if appErr != nil {
		return faceFail(c, appErr)
	}
	return c.JSON(out)
}

func faceIndexParams(c fiber.Ctx) (repository.IndexParams, *faceFault) {
	limit, bad := faceLimit(c)
	if bad != nil {
		return repository.IndexParams{}, bad
	}
	page, bad := facePage(c)
	if bad != nil {
		return repository.IndexParams{}, bad
	}
	return repository.IndexParams{
		Search: truncate(strings.TrimSpace(c.Query("q")), maxSearchLen),
		Offset: (page - 1) * limit,
		Limit:  limit,
	}, nil
}

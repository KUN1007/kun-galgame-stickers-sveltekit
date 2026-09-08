package app

import (
	"strings"

	stickerhandler "kun-galgame-sticker-api/internal/platform/sticker/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
)

// The public developer-platform face (NextMoe doc 08 §16, face name `sticker`,
// scope `sticker:read`).
//
// The gateway terminates authentication: Traefik runs ForwardAuth against
// oauth's /internal/devapi/forward-auth?face=sticker, which validates the
// nmk_ key, checks the scope, applies the key's rate and quota budget and
// meters the admitted request. Nothing reaches these routes unkeyed from the
// public internet, and the three X-NextMoe-* headers Traefik sets are the only
// trustworthy statement of who is calling. This service reads none of them --
// every answer here is the same for every caller.
//
// Traefik does not rewrite the path, so the routes must be mounted on the
// public prefix verbatim.

// facePrefixes are the public path prefixes this face answers on. Two, because
// the platform charter (§16.2, §16.5) names `/v1/<site>/*` while /v1 was
// retired wholesale on 2026-08-27 and /v2 is the only live public namespace --
// which of the two Traefik routes is the platform's call, and answering both
// costs one loop and keeps that decision from needing a redeploy here.
var facePrefixes = []string{"/v1/sticker", "/v2/sticker"}

func isFacePath(c fiber.Ctx) bool {
	path := c.Path()
	for _, prefix := range facePrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// faceCache matches catalog /v2's public lane verbatim. Every response here is
// viewer-independent, so a shared cache is correct; the cost is that a request
// Cloudflare answers is a request the platform never meters.
const faceCache = "public, max-age=300, s-maxage=1800, stale-while-revalidate=3600"

// faceCORS is deliberately not the site's CORS. The site allows one origin and
// sends credentials; the face is a public read API with no cookies, so it
// allows any origin and no credentials. Preflight is answered here for
// completeness, though a browser cannot get one past the gateway: an OPTIONS
// request carries no API key, so ForwardAuth refuses it before it arrives.
func faceCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{fiber.MethodGet, fiber.MethodOptions},
		AllowHeaders: []string{"X-API-Key", "Authorization", "Accept", "If-None-Match"},
		ExposeHeaders: []string{
			"ETag", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset",
			"X-Quota-Limit", "X-Quota-Remaining", "Retry-After",
		},
		MaxAge: 86400,
	})
}

func faceHeaders(c fiber.Ctx) error {
	c.Set(fiber.HeaderCacheControl, faceCache)
	c.Set(fiber.HeaderVary, "Accept-Encoding, Origin")
	return c.Next()
}

// mountFace registers the read face under every public prefix. There is no
// local rate limiter: the only client address this service would ever see is
// Traefik's, so an IP limiter would put every application on one budget. The
// gateway holds the real per-key budget.
func mountFace(app *fiber.App, h *stickerhandler.Handler) {
	cross := faceCORS()
	tag := etag.New()
	for _, prefix := range facePrefixes {
		face := app.Group(prefix, cross, faceHeaders, tag)

		face.Get("/packs", h.FaceListPacks)
		face.Get("/packs/:packId", h.FaceGetPack)
		face.Get("/stickers/:stickerId", h.FaceGetSticker)
		face.Get("/characters", h.FaceListCharacters)
		face.Get("/characters/:characterId", h.FaceGetCharacter)
		face.Get("/characters/:characterId/stickers", h.FaceCharacterStickers)
		face.Get("/works", h.FaceListWorks)
		face.Get("/works/:workId/packs", h.FaceWorkPacks)
		face.Get("/tags", h.FaceListTags)
	}
}

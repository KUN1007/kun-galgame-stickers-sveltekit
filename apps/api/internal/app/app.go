package app

import (
	"context"
	stderrors "errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"kun-galgame-sticker-api/internal/infrastructure/database"
	"kun-galgame-sticker-api/internal/middleware"
	identityhandler "kun-galgame-sticker-api/internal/platform/identity/handler"
	"kun-galgame-sticker-api/internal/platform/identity/oauth"
	identityservice "kun-galgame-sticker-api/internal/platform/identity/service"
	stickerhandler "kun-galgame-sticker-api/internal/platform/sticker/handler"
	stickerrepo "kun-galgame-sticker-api/internal/platform/sticker/repository"
	stickerservice "kun-galgame-sticker-api/internal/platform/sticker/service"
	"kun-galgame-sticker-api/pkg/config"
	"kun-galgame-sticker-api/pkg/errors"
	"kun-galgame-sticker-api/pkg/imageclient"
	"kun-galgame-sticker-api/pkg/response"
	"kun-galgame-sticker-api/pkg/userclient"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type App struct {
	Fiber *fiber.App
	stop  context.CancelFunc
}

// Close stops the background reference ping. Without it the goroutine outlived
// the process shutdown and kept a database handle open.
func (a *App) Close() { a.stop() }

func New(cfg *config.Config) *App {
	db := database.NewPostgres(cfg.Database, cfg.Server.Mode)

	oauthClient := oauth.NewClient(cfg.OAuth)
	authSvc := identityservice.New(oauthClient)
	authHandler := identityhandler.New(authSvc, cfg.Server.Secure)

	var imgCli *imageclient.Client
	if cfg.Image.BaseURL != "" {
		imgCli = imageclient.New(imageclient.Config{
			BaseURL:      cfg.Image.BaseURL,
			CDNBase:      cfg.Image.CDNBase,
			ClientID:     cfg.OAuth.ClientID,
			ClientSecret: cfg.OAuth.ClientSecret,
		})
	}

	users := userclient.New(userclient.Config{
		BaseURL:      cfg.OAuth.ServerURL,
		ClientID:     cfg.OAuth.ClientID,
		ClientSecret: cfg.OAuth.ClientSecret,
		ImageCDNBase: cfg.Image.CDNBase,
	})

	stickerSvc := stickerservice.New(
		stickerrepo.NewPackRepo(db),
		stickerrepo.NewStickerRepo(db),
		stickerrepo.NewTagRepo(db),
		imgCli,
		users,
	)
	ctx, stop := context.WithCancel(context.Background())
	stickerSvc.StartRefPing(ctx)
	h := stickerhandler.New(stickerSvc)

	fiberApp := fiber.New(fiber.Config{
		AppName:        "kun-galgame-sticker-api",
		ErrorHandler:   errorHandler,
		ReadBufferSize: 16 << 10,
		BodyLimit:      12 << 20,
	})
	fiberApp.Use(recover.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.CORS.AllowOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	fiberApp.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	// Auth middleware is attached per route, never with Group(prefix, mw).
	// fiber.Group registers its handlers as methodUse on the prefix, so they
	// run for every route under it regardless of which Router registered them:
	// an OptionalAuth group followed by a RequireAuth group made every write
	// resolve the session twice, and on an expired access token that meant two
	// refreshes -- the second one presenting a token the OP had already
	// rotated, which logged the user out mid-write.
	optionalAuth := middleware.OptionalAuth(authSvc, cfg.Server.Secure)
	requireAuth := middleware.RequireAuth(authSvc, cfg.Server.Secure)
	readLimit := publicLimiter()
	writeLimit := userLimiter(120)
	uploadLimit := userLimiter(60)
	cacheable := publicCache(60)

	api := fiberApp.Group("/api/v1")

	// Fully public, viewer-independent, and therefore cacheable.
	api.Get("/packs", readLimit, cacheable, etag.New(), h.ListPacks)
	api.Get("/tags", readLimit, cacheable, etag.New(), h.ListTags)

	// Public, but an author also sees their own drafts here.
	api.Get("/packs/:packId", readLimit, optionalAuth, h.GetPack)
	api.Get("/packs/:packId/download", readLimit, optionalAuth, h.DownloadPack)
	api.Get("/stickers/:stickerId", readLimit, optionalAuth, h.GetSticker)
	api.Get("/stickers/:stickerId/download", readLimit, optionalAuth, h.DownloadSticker)
	api.Get("/users/:uid/packs", readLimit, optionalAuth, h.ListUserPacks)

	api.Post("/auth/oauth/callback", writeLimit, authHandler.Callback)
	api.Post("/auth/logout", authHandler.Logout)
	api.Get("/auth/me", optionalAuth, authHandler.Me)

	api.Get("/me/packs", requireAuth, h.ListMyPacks)
	api.Post("/me/packs", requireAuth, writeLimit, h.CreatePack)
	api.Patch("/me/packs/:packId", requireAuth, writeLimit, h.PatchPack)
	api.Delete("/me/packs/:packId", requireAuth, writeLimit, h.DeletePack)
	api.Post("/me/packs/:packId/publish", requireAuth, writeLimit, h.Publish)
	api.Post("/me/packs/:packId/unpublish", requireAuth, writeLimit, h.Unpublish)
	api.Post("/me/packs/:packId/images", requireAuth, uploadLimit, h.UploadImage)
	api.Post("/me/packs/:packId/stickers", requireAuth, writeLimit, h.AddSticker)
	api.Patch("/me/packs/:packId/stickers/:stickerId", requireAuth, writeLimit, h.PatchSticker)
	api.Delete("/me/packs/:packId/stickers/:stickerId", requireAuth, writeLimit, h.DeleteSticker)
	api.Put("/me/packs/:packId/stickers/order", requireAuth, writeLimit, h.ReorderStickers)

	return &App{Fiber: fiberApp, stop: stop}
}

func publicLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:          300,
		Expiration:   time.Minute,
		LimitReached: tooManyRequests,
	})
}

// userLimiter counts per signed-in user, falling back to the client address for
// anonymous callers -- a whole office behind one NAT would otherwise share a
// single write budget.
func userLimiter(perMinute int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        perMinute,
		Expiration: time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			if user := middleware.CurrentUser(c); user != nil {
				return "u" + strconv.Itoa(user.ID)
			}
			return c.IP()
		},
		LimitReached: tooManyRequests,
	})
}

func tooManyRequests(c fiber.Ctx) error {
	return response.Error(c, errors.New(errors.CodeBiz, "too many requests", 429))
}

func publicCache(seconds int) fiber.Handler {
	value := "public, max-age=" + strconv.Itoa(seconds)
	return func(c fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, value)
		return c.Next()
	}
}

func errorHandler(c fiber.Ctx, err error) error {
	var appErr *errors.AppError
	if stderrors.As(err, &appErr) {
		return response.Error(c, appErr)
	}
	var fe *fiber.Error
	if stderrors.As(err, &fe) {
		return response.Error(c, errors.New(errors.CodeBiz, fe.Message, fe.Code))
	}
	slog.Error("unhandled", "error", err)
	return response.Error(c, errors.ErrInternal("internal error"))
}

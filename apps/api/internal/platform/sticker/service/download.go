package service

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"kun-galgame-sticker-api/internal/platform/sticker/model"
	"kun-galgame-sticker-api/pkg/errors"

	"github.com/google/uuid"
)

const downloadTimeout = 60 * time.Second

// Downloads are proxied rather than handed to the browser as CDN links: a
// cross-origin <a download> is ignored by browsers and a fetch+blob needs CORS
// on the image host. Going through the API also makes download_count real.

func (s *Service) PackForDownload(id uuid.UUID, v Viewer) (*model.Pack, []model.Sticker, *errors.AppError) {
	pack, appErr := s.visiblePack(id, v)
	if appErr != nil {
		return nil, nil, appErr
	}
	rows, err := s.stickers.ListByPack(pack.ID)
	if err != nil {
		return nil, nil, errors.ErrInternal("failed to load stickers")
	}
	if len(rows) == 0 {
		return nil, nil, errors.ErrPackNeedsSticker()
	}
	return pack, rows, nil
}

func (s *Service) StickerForDownload(id uuid.UUID, v Viewer) (*model.Sticker, *errors.AppError) {
	return s.visibleSticker(id, v)
}

func (s *Service) BumpDownload(packID uuid.UUID) { _ = s.packs.BumpDownload(packID) }

func (s *Service) WriteZip(ctx context.Context, w io.Writer, stickers []model.Sticker) error {
	ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
	defer cancel()

	archive := zip.NewWriter(w)
	for _, row := range stickers {
		body, _, err := s.openImage(ctx, row.ImageHash)
		if err != nil {
			archive.Close()
			return err
		}
		entry, createErr := archive.Create(fmt.Sprintf("%03d.webp", row.Position))
		if createErr != nil {
			body.Close()
			archive.Close()
			return createErr
		}
		_, copyErr := io.Copy(entry, body)
		body.Close()
		if copyErr != nil {
			archive.Close()
			return copyErr
		}
	}
	return archive.Close()
}

func (s *Service) OpenImage(ctx context.Context, hash string) (io.ReadCloser, int64, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
	body, size, err := s.openImage(ctx, hash)
	if err != nil {
		cancel()
		return nil, 0, errors.ErrImageUnavailable()
	}
	return &cancelOnClose{ReadCloser: body, cancel: cancel}, size, nil
}

func (s *Service) openImage(ctx context.Context, hash string) (io.ReadCloser, int64, error) {
	main, _ := s.urls(hash)
	if main == "" {
		return nil, 0, fmt.Errorf("no image url for %s", hash)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, main, nil)
	if err != nil {
		return nil, 0, err
	}
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("image cdn returned %d for %s", resp.StatusCode, hash)
	}
	return resp.Body, resp.ContentLength, nil
}

type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c *cancelOnClose) Close() error {
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

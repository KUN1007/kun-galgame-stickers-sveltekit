package service

import (
	"context"
	"log/slog"
	"time"
)

const (
	refPingInterval = 24 * time.Hour
	refPingDelay    = time.Minute
	refPingBatch    = 1000
)

// StartRefPing keeps every image this site still references warm. The image
// service reports images cold after 60 days without a ping and soft-deletes at
// 365, so missing a day is survivable but missing months is not.
func (s *Service) StartRefPing(ctx context.Context) {
	if s.images == nil {
		return
	}
	go func() {
		timer := time.NewTimer(refPingDelay)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}

		ticker := time.NewTicker(refPingInterval)
		defer ticker.Stop()
		for {
			s.PingHashes(ctx)
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

func (s *Service) PingHashes(ctx context.Context) {
	hashes, err := s.stickers.ListImageHashes()
	if err != nil {
		slog.Error("refping: list hashes", "error", err)
		return
	}
	if len(hashes) == 0 {
		return
	}

	var pinged, failed int
	for start := 0; start < len(hashes); start += refPingBatch {
		end := min(start+refPingBatch, len(hashes))
		// One bad batch used to abort the whole sweep, so a single transient
		// 5xx silently skipped that day's refresh for every image.
		result, err := s.images.ReferencePing(ctx, hashes[start:end])
		if err != nil {
			failed += end - start
			slog.Warn("refping: batch failed", "from", start, "to", end, "error", err)
			continue
		}
		pinged += end - start
		if len(result.NotFound) > 0 {
			slog.Warn("refping: hashes unknown to the image service",
				"count", len(result.NotFound), "sample", result.NotFound[0])
		}
	}
	slog.Info("refping done", "pinged", pinged, "failed", failed)
}

package service

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"strings"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/pkg/errors"
)

// avatarVariant is the image-service variant the fallback pool is served at.
//
// KunAvatar's largest size renders at 48 CSS px, so 128 covers a 2x display
// exactly. It is also 8x cheaper than the original: measured on the pool's own
// bytes, main is 44,018 B, _320 is 22,170 B and _128 is 5,318 B. At this
// volume the variant choice is the single largest lever after caching -- the
// site's own grid keeps using _320 (thumbVariant), which would be four times
// the bytes for an avatar nobody can see the detail of.
const avatarVariant = "128"

// AvatarPool is the ecosystem-wide default-avatar manifest.
//
// It exists because 94.3% of NextMoe accounts have no avatar image, so a
// sticker is what renders in their place -- tens of millions of image requests
// a day. Every consuming site used to compute those URLs itself from a
// hardcoded path on this host, which is why they all broke at once when the
// static files went away.
//
// The contract is deliberately narrow: this returns READY, ABSOLUTE,
// content-addressed CDN URLs, never a hash plus a rule for assembling one. A
// consumer that assembles URLs has re-acquired the coupling this endpoint
// exists to remove; changing CDN host or variant would break it again. Callers
// treat the list as an opaque fixed-length array and pick with
// hash(seed) % len(urls).
//
// It is meant to be fetched SERVER-SIDE, once an hour, by each site -- never
// per render and never by a browser. A few requests a day in total.
func (s *Service) AvatarPool() (*dto.AvatarPool, *errors.AppError) {
	hashes, err := s.stickers.AvatarPoolHashes()
	if err != nil {
		slog.Error("avatar pool query failed", "error", err)
		return nil, errors.ErrInternal("failed to load the avatar pool")
	}

	urls := make([]string, 0, len(hashes))
	for _, hash := range hashes {
		if s.images == nil {
			break
		}
		urls = append(urls, s.images.VariantURL(hash, avatarVariant))
	}

	// A content-derived version: it changes when, and only when, the pool
	// does. Consumers can use it to tell "my cached copy is current" from "the
	// fetch failed and I am on the baked fallback" -- an ETag would only tell
	// them about one hop.
	sum := sha256.Sum256([]byte(strings.Join(hashes, "\n")))
	return &dto.AvatarPool{
		Version: hex.EncodeToString(sum[:])[:16],
		Variant: avatarVariant,
		URLs:    urls,
	}, nil
}

package repository

import (
	"errors"

	"gorm.io/gorm/clause"
)

var (
	errStickerLimit    = errors.New("sticker limit reached")
	errReorderMismatch = errors.New("reorder list does not match the pack")
)

func IsStickerLimit(err error) bool    { return errors.Is(err, errStickerLimit) }
func IsReorderMismatch(err error) bool { return errors.Is(err, errReorderMismatch) }

func lockForUpdate() clause.Locking { return clause.Locking{Strength: "UPDATE"} }

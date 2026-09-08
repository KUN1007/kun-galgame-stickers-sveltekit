package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	PackDraft     int16 = 0
	PackPublished int16 = 1
	PackHidden    int16 = 2
	PackRemoved   int16 = 3
)

const (
	RatingSFW  int16 = 0
	RatingNSFW int16 = 1
)

type Pack struct {
	ID             uuid.UUID      `gorm:"column:id;primaryKey;default:uuidv7()"`
	OwnerUID       int            `gorm:"column:owner_uid"`
	Status         int16          `gorm:"column:status"`
	IsOfficial     bool           `gorm:"column:is_official"`
	ContentRating  int16          `gorm:"column:content_rating"`
	Title          datatypes.JSON `gorm:"column:title;type:jsonb"`
	Description    datatypes.JSON `gorm:"column:description;type:jsonb"`
	CoverStickerID *uuid.UUID     `gorm:"column:cover_sticker_id"`
	StickerCount   int            `gorm:"column:sticker_count"`
	ViewCount      int64          `gorm:"column:view_count"`
	DownloadCount  int64          `gorm:"column:download_count"`
	SearchText     string         `gorm:"column:search_text"`
	// CatalogWorkID is the infra catalog id of the game this pack is about,
	// when its author declared one. The name and cover beside it are a display
	// snapshot so a list page never calls catalog; catalog stays the source of
	// truth for the identity itself.
	CatalogWorkID    *int64         `gorm:"column:catalog_work_id"`
	CatalogWorkName  datatypes.JSON `gorm:"column:catalog_work_name;type:jsonb"`
	CatalogWorkCover string         `gorm:"column:catalog_work_cover"`
	CreatedAt        time.Time      `gorm:"column:created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at"`
	PublishedAt      *time.Time     `gorm:"column:published_at"`
}

func (Pack) TableName() string { return "pack" }

// The jsonb columns are NOT NULL DEFAULT '{}', but GORM writes an explicit NULL
// for a nil datatypes.JSON instead of omitting the column, so a caller that
// leaves one unset gets a constraint violation rather than the default.
func (p *Pack) BeforeSave(*gorm.DB) error {
	p.Title = orEmptyJSON(p.Title)
	p.Description = orEmptyJSON(p.Description)
	p.CatalogWorkName = orEmptyJSON(p.CatalogWorkName)
	return nil
}

func orEmptyJSON(in datatypes.JSON) datatypes.JSON {
	if len(in) == 0 {
		return datatypes.JSON("{}")
	}
	return in
}

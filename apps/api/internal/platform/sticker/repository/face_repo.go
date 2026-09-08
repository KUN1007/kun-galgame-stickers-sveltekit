package repository

import (
	"kun-galgame-sticker-api/internal/platform/sticker/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// The queries here back the public face's two catalog-id entry points -- "all
// stickers of this character", "all packs about this game" -- plus the indexes
// that let a caller discover which ids this site actually holds. They differ
// from the site's own queries in one way that matters: they are unconditionally
// scoped to published packs, because there is no viewer to widen them.

// IndexParams pages one of the catalog-identity indexes.
type IndexParams struct {
	Search string
	// Work narrows the character index to one game. Ignored by WorkIndex.
	Work   int64
	Offset int
	Limit  int
}

// CharacterRow is one catalog character this site has stickers of.
type CharacterRow struct {
	CatalogCharacterID    int64
	CatalogCharacterName  datatypes.JSON
	CatalogCharacterImage string
	CatalogWorkID         *int64
	CatalogWorkName       datatypes.JSON
	StickerCount          int64
}

// WorkRow is one catalog work this site has material for.
type WorkRow struct {
	CatalogWorkID     int64
	CatalogWorkName   datatypes.JSON
	CatalogWorkCover  string
	CatalogWorkRating string
	StickerCount      int64
}

// characterSelect is shared by the index and the single lookup so the two can
// never disagree about what a character row contains.
const characterSelect = `sticker.catalog_character_id,
	min(sticker.catalog_character_name::text)::jsonb AS catalog_character_name,
	min(sticker.catalog_character_image) AS catalog_character_image,
	min(sticker.catalog_work_id) AS catalog_work_id,
	min(sticker.catalog_work_name::text)::jsonb AS catalog_work_name,
	count(*) AS sticker_count`

// publishedStickers is the base every face query starts from. The join is what
// keeps a draft's stickers out of the public index: sticker rows themselves
// carry no visibility of their own.
func (r *StickerRepo) publishedStickers() *gorm.DB {
	return r.db.Model(&model.Sticker{}).
		Joins("JOIN pack ON pack.id = sticker.pack_id AND pack.status = ?", model.PackPublished)
}

// CharacterIndex lists the catalog characters this site has stickers of, most
// stickers first. min() picks an arbitrary row's snapshot per character, which
// is fine: every row for one catalog id was written from the same catalog
// identity.
func (r *StickerRepo) CharacterIndex(p IndexParams) ([]CharacterRow, int64, error) {
	base := func() *gorm.DB {
		q := r.publishedStickers().Where("sticker.catalog_character_id IS NOT NULL")
		if p.Search != "" {
			q = q.Where("sticker.catalog_character_name::text ILIKE ?", "%"+escapeLike(p.Search)+"%")
		}
		if p.Work > 0 {
			q = q.Where("sticker.catalog_work_id = ?", p.Work)
		}
		return q
	}

	var total int64
	if err := base().Distinct("sticker.catalog_character_id").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	var rows []CharacterRow
	err := base().
		Select(characterSelect).
		Group("sticker.catalog_character_id").
		Order("count(*) DESC, sticker.catalog_character_id ASC").
		Offset(p.Offset).Limit(p.Limit).
		Scan(&rows).Error
	return rows, total, err
}

// Character loads one row of the index. It answers nil rather than an error
// when the id is unknown here -- catalog may well know the character, this
// site simply has no stickers of them.
func (r *StickerRepo) Character(characterID int64) (*CharacterRow, error) {
	var row CharacterRow
	err := r.publishedStickers().
		Where("sticker.catalog_character_id = ?", characterID).
		Select(characterSelect).
		Group("sticker.catalog_character_id").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.CatalogCharacterID == 0 {
		return nil, nil
	}
	return &row, nil
}

// WorkIndex lists the catalog works this site has material for. It counts
// stickers only: how many packs a game is "in" is answered by
// /works/{id}/packs, and having the number in two places invites the two to
// disagree -- a pack can declare a game at pack level while none of its
// stickers name it, and vice versa.
func (r *StickerRepo) WorkIndex(p IndexParams) ([]WorkRow, int64, error) {
	base := func() *gorm.DB {
		q := r.publishedStickers().Where("sticker.catalog_work_id IS NOT NULL")
		if p.Search != "" {
			q = q.Where("sticker.catalog_work_name::text ILIKE ?", "%"+escapeLike(p.Search)+"%")
		}
		return q
	}

	var total int64
	if err := base().Distinct("sticker.catalog_work_id").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	var rows []WorkRow
	err := base().
		Select(`sticker.catalog_work_id,
			min(sticker.catalog_work_name::text)::jsonb AS catalog_work_name,
			min(sticker.catalog_work_rating) AS catalog_work_rating,
			min(pack.catalog_work_cover) AS catalog_work_cover,
			count(*) AS sticker_count`).
		Group("sticker.catalog_work_id").
		Order("count(*) DESC, sticker.catalog_work_id ASC").
		Offset(p.Offset).Limit(p.Limit).
		Scan(&rows).Error
	return rows, total, err
}

// StickersByCharacter pages one character's stickers, newest first, published
// packs only.
func (r *StickerRepo) StickersByCharacter(characterID int64, offset, limit int) ([]model.Sticker, int64, error) {
	var total int64
	if err := r.publishedStickers().
		Where("sticker.catalog_character_id = ?", characterID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}
	var rows []model.Sticker
	err := r.publishedStickers().
		Where("sticker.catalog_character_id = ?", characterID).
		Select("sticker.*").
		Order("sticker.created_at DESC, sticker.id DESC").
		Offset(offset).Limit(limit).
		Find(&rows).Error
	return rows, total, err
}

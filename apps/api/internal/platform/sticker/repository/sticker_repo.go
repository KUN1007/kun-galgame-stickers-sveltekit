package repository

import (
	"time"

	"kun-galgame-sticker-api/internal/platform/sticker/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StickerRepo struct{ db *gorm.DB }

func NewStickerRepo(db *gorm.DB) *StickerRepo { return &StickerRepo{db: db} }

func (r *StickerRepo) ListByPack(packID uuid.UUID) ([]model.Sticker, error) {
	var rows []model.Sticker
	err := r.db.Where("pack_id = ?", packID).Order("position ASC").Find(&rows).Error
	return rows, err
}

func (r *StickerRepo) Get(id uuid.UUID) (*model.Sticker, error) {
	var row model.Sticker
	if err := r.db.First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ByIDs backs the cover lookup for a page of packs: one query for the whole
// page instead of one per pack.
func (r *StickerRepo) ByIDs(ids []uuid.UUID) (map[uuid.UUID]model.Sticker, error) {
	out := make(map[uuid.UUID]model.Sticker, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var rows []model.Sticker
	if err := r.db.Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = row
	}
	return out, nil
}

func (r *StickerRepo) CountByPack(packID uuid.UUID) (int64, error) {
	var n int64
	err := r.db.Model(&model.Sticker{}).Where("pack_id = ?", packID).Count(&n).Error
	return n, err
}

func (r *StickerRepo) CountRecentByOwner(uid int, since time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&model.Sticker{}).
		Joins("JOIN pack ON pack.id = sticker.pack_id").
		Where("pack.owner_uid = ? AND sticker.created_at >= ?", uid, since).
		Count(&n).Error
	return n, err
}

// Append allocates the next position under a row lock on the owning pack.
// Without the lock two concurrent uploads both read the same max(position) and
// the second insert dies on sticker_pack_position_idx.
func (r *StickerRepo) Append(packID uuid.UUID, row *model.Sticker, limit int) (int64, error) {
	var count int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var locked model.Pack
		if err := tx.Clauses(lockForUpdate()).First(&locked, "id = ?", packID).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Sticker{}).Where("pack_id = ?", packID).Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(limit) {
			return errStickerLimit
		}

		var maxPos *int
		if err := tx.Model(&model.Sticker{}).Where("pack_id = ?", packID).
			Select("max(position)").Scan(&maxPos).Error; err != nil {
			return err
		}
		row.PackID = packID
		row.Position = 1
		if maxPos != nil {
			row.Position = *maxPos + 1
		}
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		count++
		return NewPackRepo(tx).SyncStickerCount(tx, packID)
	})
	return count, err
}

func (r *StickerRepo) Save(row *model.Sticker) error { return r.db.Save(row).Error }

func (r *StickerRepo) Delete(packID, id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&model.Sticker{}, "id = ? AND pack_id = ?", id, packID)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return NewPackRepo(tx).SyncStickerCount(tx, packID)
	})
}

// Reorder rewrites positions to match the given order. Positions are shifted
// out of range first: they are unique per pack, so assigning 1..n directly
// collides with the rows that still hold those numbers.
func (r *StickerRepo) Reorder(packID uuid.UUID, ids []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var locked model.Pack
		if err := tx.Clauses(lockForUpdate()).First(&locked, "id = ?", packID).Error; err != nil {
			return err
		}
		var existing []model.Sticker
		if err := tx.Where("pack_id = ?", packID).Find(&existing).Error; err != nil {
			return err
		}
		if len(existing) != len(ids) {
			return errReorderMismatch
		}
		known := make(map[uuid.UUID]bool, len(existing))
		for _, row := range existing {
			known[row.ID] = true
		}
		for _, id := range ids {
			if !known[id] {
				return errReorderMismatch
			}
			delete(known, id)
		}

		if err := tx.Model(&model.Sticker{}).Where("pack_id = ?", packID).
			UpdateColumn("position", gorm.Expr("-position")).Error; err != nil {
			return err
		}
		for i, id := range ids {
			if err := tx.Model(&model.Sticker{}).
				Where("id = ? AND pack_id = ?", id, packID).
				UpdateColumn("position", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *StickerRepo) ListImageHashes() ([]string, error) {
	var hashes []string
	err := r.db.Model(&model.Sticker{}).
		Where("image_hash <> ''").
		Distinct().Pluck("image_hash", &hashes).Error
	return hashes, err
}

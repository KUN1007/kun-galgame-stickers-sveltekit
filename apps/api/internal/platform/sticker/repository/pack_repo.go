package repository

import (
	"strings"

	"kun-galgame-sticker-api/internal/platform/sticker/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListParams struct {
	OwnerUID     int
	Statuses     []int16
	OfficialOnly bool
	SFWOnly      bool
	Search       string
	TagSlug      string
	LinkedOnly   bool
	CatalogWork  int64
	// AnyCatalogWork matches a pack that declares the work or holds a sticker
	// of it.
	AnyCatalogWork int64
	Order          string
	Offset         int
	Limit          int
}

type PackRepo struct{ db *gorm.DB }

func NewPackRepo(db *gorm.DB) *PackRepo { return &PackRepo{db: db} }

func (r *PackRepo) DB() *gorm.DB { return r.db }

func (r *PackRepo) filtered(p ListParams) *gorm.DB {
	q := r.db.Model(&model.Pack{})
	if p.OwnerUID > 0 {
		q = q.Where("owner_uid = ?", p.OwnerUID)
	}
	if len(p.Statuses) > 0 {
		q = q.Where("status IN ?", p.Statuses)
	}
	if p.OfficialOnly {
		q = q.Where("is_official")
	}
	if p.SFWOnly {
		q = q.Where("content_rating = ?", model.RatingSFW)
	}
	if p.Search != "" {
		q = q.Where("search_text ILIKE ?", "%"+escapeLike(p.Search)+"%")
	}
	if p.LinkedOnly {
		q = q.Where("catalog_work_id IS NOT NULL")
	}
	if p.CatalogWork > 0 {
		q = q.Where("catalog_work_id = ?", p.CatalogWork)
	}
	if p.AnyCatalogWork > 0 {
		q = q.Where(
			"(catalog_work_id = ? OR EXISTS (SELECT 1 FROM sticker s WHERE s.pack_id = pack.id AND s.catalog_work_id = ?))",
			p.AnyCatalogWork, p.AnyCatalogWork,
		)
	}
	if p.TagSlug != "" {
		q = q.Where(
			"EXISTS (SELECT 1 FROM pack_tag pt JOIN tag t ON t.id = pt.tag_id WHERE pt.pack_id = pack.id AND t.slug = ?)",
			p.TagSlug,
		)
	}
	return q
}

func (r *PackRepo) List(p ListParams) ([]model.Pack, int64, error) {
	var total int64
	if err := r.filtered(p).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}
	var rows []model.Pack
	err := r.filtered(p).Order(p.Order).Offset(p.Offset).Limit(p.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *PackRepo) Get(id uuid.UUID) (*model.Pack, error) {
	var row model.Pack
	if err := r.db.First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ByIDs loads a set of packs in one query. Callers that start from stickers
// (the character page) need this to avoid a query per sticker.
func (r *PackRepo) ByIDs(ids []string) ([]model.Pack, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []model.Pack
	err := r.db.Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *PackRepo) CountByOwner(uid int) (int64, error) {
	var n int64
	err := r.db.Model(&model.Pack{}).Where("owner_uid = ?", uid).Count(&n).Error
	return n, err
}

func (r *PackRepo) Create(row *model.Pack) error { return r.db.Create(row).Error }

func (r *PackRepo) Save(row *model.Pack) error { return r.db.Save(row).Error }

func (r *PackRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Pack{}, "id = ?", id).Error
}

func (r *PackRepo) BumpDownload(id uuid.UUID) error {
	return r.db.Model(&model.Pack{}).Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

func (r *PackRepo) BumpView(id uuid.UUID) error {
	return r.db.Model(&model.Pack{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// SyncStickerCount recomputes rather than increments: the callers that change
// the count (add, delete, reorder-with-delete) can run concurrently, and a
// recount from the same transaction cannot drift.
func (r *PackRepo) SyncStickerCount(tx *gorm.DB, packID uuid.UUID) error {
	return tx.Model(&model.Pack{}).Where("id = ?", packID).
		UpdateColumn("sticker_count",
			tx.Model(&model.Sticker{}).Select("count(*)").Where("pack_id = ?", packID),
		).Error
}

func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}

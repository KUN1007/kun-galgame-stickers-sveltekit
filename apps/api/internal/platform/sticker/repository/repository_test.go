package repository_test

import (
	"os"
	"strings"
	"sync"
	"testing"

	"kun-galgame-sticker-api/internal/platform/sticker/model"
	"kun-galgame-sticker-api/internal/platform/sticker/repository"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// These tests need a throwaway database with the migrations applied:
//
//	createdb kungalgame_sticker_test
//	KUN_DATABASE_URL=<test dsn> go run ./cmd/migrate -dir up
//	STICKER_TEST_DSN=<test dsn> go test ./internal/platform/sticker/repository/
//
// Without STICKER_TEST_DSN they skip, so `go test ./...` stays green on a
// machine with no database.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("STICKER_TEST_DSN")
	if dsn == "" {
		t.Skip("set STICKER_TEST_DSN to run repository tests")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Discard,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	return db
}

func newPack(t *testing.T, db *gorm.DB) *model.Pack {
	t.Helper()
	pack := &model.Pack{
		OwnerUID: 4242,
		Status:   model.PackDraft,
		Title:    []byte(`{"zh-cn":"test"}`),
	}
	if err := repository.NewPackRepo(db).Create(pack); err != nil {
		t.Fatalf("create pack: %v", err)
	}
	t.Cleanup(func() { db.Delete(&model.Pack{}, "id = ?", pack.ID) })
	return pack
}

func sticker() *model.Sticker {
	return &model.Sticker{ImageHash: strings.Repeat("a", 64)}
}

// Positions are unique per pack and used to be handed out by max(position)+1
// with no lock, so two simultaneous uploads picked the same number and the
// second insert died on the unique index.
func TestAppendIsSafeUnderConcurrency(t *testing.T) {
	db := testDB(t)
	pack := newPack(t, db)
	repo := repository.NewStickerRepo(db)

	const workers = 12
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			if _, err := repo.Append(pack.ID, sticker(), 100); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent append: %v", err)
	}

	rows, err := repo.ListByPack(pack.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rows) != workers {
		t.Fatalf("stored %d stickers, want %d", len(rows), workers)
	}
	for i, row := range rows {
		if row.Position != i+1 {
			t.Fatalf("position %d at index %d, want %d", row.Position, i, i+1)
		}
	}

	fresh, err := repository.NewPackRepo(db).Get(pack.ID)
	if err != nil {
		t.Fatalf("reload pack: %v", err)
	}
	if fresh.StickerCount != workers {
		t.Fatalf("sticker_count = %d, want %d", fresh.StickerCount, workers)
	}
}

func TestAppendRefusesPastTheLimit(t *testing.T) {
	db := testDB(t)
	pack := newPack(t, db)
	repo := repository.NewStickerRepo(db)

	for range 3 {
		if _, err := repo.Append(pack.ID, sticker(), 3); err != nil {
			t.Fatalf("append: %v", err)
		}
	}
	_, err := repo.Append(pack.ID, sticker(), 3)
	if !repository.IsStickerLimit(err) {
		t.Fatalf("append past the limit returned %v, want the limit error", err)
	}
}

func TestReorder(t *testing.T) {
	db := testDB(t)
	pack := newPack(t, db)
	repo := repository.NewStickerRepo(db)

	var ids []uuid.UUID
	for range 3 {
		row := sticker()
		if _, err := repo.Append(pack.ID, row, 10); err != nil {
			t.Fatalf("append: %v", err)
		}
		ids = append(ids, row.ID)
	}

	reversed := []uuid.UUID{ids[2], ids[1], ids[0]}
	if err := repo.Reorder(pack.ID, reversed); err != nil {
		t.Fatalf("reorder: %v", err)
	}
	rows, err := repo.ListByPack(pack.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for i, row := range rows {
		if row.ID != reversed[i] {
			t.Fatalf("position %d holds %s, want %s", i+1, row.ID, reversed[i])
		}
	}

	if err := repo.Reorder(pack.ID, ids[:2]); !repository.IsReorderMismatch(err) {
		t.Fatalf("short reorder list returned %v, want a mismatch error", err)
	}
	if err := repo.Reorder(pack.ID, []uuid.UUID{ids[0], ids[1], uuid.New()}); !repository.IsReorderMismatch(err) {
		t.Fatalf("reorder with a foreign id returned %v, want a mismatch error", err)
	}
}

func TestDeleteKeepsTheCountInSync(t *testing.T) {
	db := testDB(t)
	pack := newPack(t, db)
	repo := repository.NewStickerRepo(db)

	row := sticker()
	if _, err := repo.Append(pack.ID, row, 10); err != nil {
		t.Fatalf("append: %v", err)
	}
	if err := repo.Delete(pack.ID, row.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	fresh, err := repository.NewPackRepo(db).Get(pack.ID)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if fresh.StickerCount != 0 {
		t.Fatalf("sticker_count = %d after delete, want 0", fresh.StickerCount)
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Cat Girl":    "cat-girl",
		"  cat girl ": "cat-girl",
		"猫娘":          "猫娘",
		"日常 / 吐槽":     "日常-吐槽",
		"!!!":         "",
		"a---b":       "a-b",
	}
	for in, want := range cases {
		if got := repository.Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

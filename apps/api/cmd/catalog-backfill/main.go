// Command catalog-backfill attaches catalog identities to the stickers that
// were annotated by hand before catalog existed.
//
// The seven seeded official packs carry a game and a character name per
// sticker as free text, in romaji (en-us) and Chinese (zh-cn). Those names are
// good data -- they were transcribed from the games -- they simply point at
// nothing. This resolves each one against catalog and writes the id plus the
// same display snapshot the API writes when an author picks a game by hand.
//
// It refuses to guess. A game or a character is only written when a candidate
// matches on a normalized name; everything else is printed for a person to
// resolve, because a wrong identity is worse than a missing one -- it puts a
// sticker on some other character's page and there is no signal it is wrong.
//
//	go run ./cmd/catalog-backfill              # dry run, prints the plan
//	go run ./cmd/catalog-backfill -apply       # writes
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"kun-galgame-sticker-api/internal/infrastructure/database"
	"kun-galgame-sticker-api/internal/platform/sticker/model"
	"kun-galgame-sticker-api/pkg/catalogclient"
	"kun-galgame-sticker-api/pkg/config"

	"github.com/joho/godotenv"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func main() {
	apply := flag.Bool("apply", false, "write the resolved links (default: dry run)")
	refresh := flag.Bool("refresh", false, "also re-resolve stickers that already carry a link")
	packLimit := flag.Int("packs", 0, "only touch this many packs (0 = all)")
	flag.Parse()

	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	catalog := catalogclient.New(catalogclient.Config{
		BaseURL: cfg.Catalog.BaseURL,
		APIKey:  cfg.Catalog.APIKey,
	})
	if !catalog.Configured() {
		log.Fatal("catalog is not configured: set KUN_CATALOG_BASE_URL and KUN_CATALOG_API_KEY")
	}

	db := database.NewPostgres(cfg.Database, "prod")
	run := &backfill{db: db, catalog: catalog, apply: *apply, refresh: *refresh, works: map[string]*match{}}
	if err := run.do(context.Background(), *packLimit); err != nil {
		log.Fatalf("backfill: %v", err)
	}
}

// Identities a person resolved by reading catalog, for the cases no exact-name
// rule can reach. Keyed by the transcription as it appears on the stickers.
var (
	// The stickers say "Summer Pockets"; catalog's entry carries the original's
	// release date (2018-04-25) and the plain Chinese title 夏日口袋 while its
	// display name has drifted to the REFLECTION BLUE edition. It is the same
	// identity the other 497 stickers resolved to, and its roster holds the
	// character on this one (Tsumugi Wenders).
	manualWorks = map[string]int64{
		"summerpockets": 4,
	}

	// "Tina" is how the stickers name ティナ・フルール・レニスフィア. A short
	// form cannot match by rule without opening the door to matching any
	// character whose name merely starts the same way.
	manualCharacters = map[string]int64{
		"2751|tina": 44629,
	}
)

type match struct {
	work    *catalogclient.Work
	roster  []catalogclient.Character
	missing bool
}

type backfill struct {
	db      *gorm.DB
	catalog *catalogclient.Client
	apply   bool
	// refresh re-resolves rows that already carry a link, for when the snapshot
	// gains a field or catalog renames something. Off by default: the ordinary
	// run should only ever add.
	refresh bool

	works map[string]*match // keyed by the sticker's romaji game name

	linked         int
	charLinked     int
	unknownWork    map[string]int
	unknownChar    map[string]int
	touchedPackIDs map[string]bool
}

func (b *backfill) do(ctx context.Context, packLimit int) error {
	b.unknownWork = map[string]int{}
	b.unknownChar = map[string]int{}
	b.touchedPackIDs = map[string]bool{}

	var rows []model.Sticker
	q := b.db.Where("game::text <> '{}'").Order("pack_id, position")
	if !b.refresh {
		q = q.Where("catalog_work_id IS NULL")
	}
	if packLimit > 0 {
		var packIDs []string
		if err := b.db.Model(&model.Pack{}).Order("created_at").Limit(packLimit).
			Pluck("id::text", &packIDs).Error; err != nil {
			return err
		}
		q = q.Where("pack_id IN ?", packIDs)
	}
	if err := q.Find(&rows).Error; err != nil {
		return err
	}
	fmt.Printf("%d stickers to resolve\n\n", len(rows))

	for i := range rows {
		if err := b.one(ctx, &rows[i]); err != nil {
			return err
		}
	}

	if b.apply {
		for packID := range b.touchedPackIDs {
			if err := b.syncSearchText(packID); err != nil {
				return err
			}
		}
	}
	b.report()
	return nil
}

func (b *backfill) one(ctx context.Context, row *model.Sticker) error {
	game := decode(row.Game)
	character := decode(row.CharacterName)
	romaji, chinese := game["en-us"], game["zh-cn"]
	if romaji == "" && chinese == "" {
		return nil
	}

	found, err := b.resolveWork(ctx, romaji, chinese)
	if err != nil {
		return err
	}
	if found.missing {
		b.unknownWork[label(romaji, chinese)]++
		return nil
	}

	workID := mustID(found.work.ID)
	workName := fillMissing(
		catalogclient.LocalizedName(found.work.DisplayName, found.work.Localized, found.work.Latin), game)

	var charID *int64
	charName := map[string]string{}
	charImage := ""
	hit := pickCharacter(found.roster, character["en-us"], character["zh-cn"])
	if hit == nil {
		if id, ok := manualCharacters[fmt.Sprintf("%d|%s", workID, tight(character["en-us"]))]; ok {
			hit = findInRoster(found.roster, id)
		}
	}
	if hit != nil {
		id := mustID(hit.ID)
		charID = &id
		charName = fillMissing(
			catalogclient.LocalizedName(hit.DisplayName, hit.Localized, hit.Latin), character)
		if hit.Image != nil {
			charImage = hit.Image.URL
		}
		b.charLinked++
	} else if character["en-us"] != "" || character["zh-cn"] != "" {
		b.unknownChar[label(romaji, chinese)+" :: "+label(character["en-us"], character["zh-cn"])]++
	}
	b.linked++
	b.touchedPackIDs[row.PackID.String()] = true

	if !b.apply {
		return nil
	}
	return b.db.Model(&model.Sticker{}).Where("id = ?", row.ID).Updates(map[string]any{
		"catalog_work_id":         workID,
		"catalog_work_name":       encode(workName),
		"catalog_work_rating":     found.work.ContentRating,
		"catalog_character_id":    charID,
		"catalog_character_name":  encode(charName),
		"catalog_character_image": charImage,
		"updated_at":              time.Now(),
	}).Error
}

// resolveWork searches on both transcriptions and keeps the result, hit or
// miss, so ninety-seven games cost ninety-seven lookups rather than one per
// sticker.
func (b *backfill) resolveWork(ctx context.Context, romaji, chinese string) (*match, error) {
	key := norm(romaji) + "|" + norm(chinese)
	if cached, ok := b.works[key]; ok {
		return cached, nil
	}
	found := &match{missing: true}
	b.works[key] = found

	var candidates []catalogclient.Work
	for _, query := range []string{romaji, chinese} {
		if strings.TrimSpace(query) == "" {
			continue
		}
		hits, err := b.catalog.SearchWorks(ctx, query, 20)
		if err != nil {
			return nil, fmt.Errorf("search %q: %w", query, err)
		}
		candidates = append(candidates, hits...)
	}

	best := pickWork(candidates, romaji, chinese)
	if best == nil {
		id, ok := manualWorks[tight(romaji)]
		if !ok {
			return found, nil
		}
		resolved, err := b.catalog.Work(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("manual work %d: %w", id, err)
		}
		best = resolved
	}
	roster, err := b.catalog.WorkCharacters(ctx, mustID(best.ID))
	if err != nil {
		return nil, fmt.Errorf("roster for %s: %w", best.DisplayName, err)
	}
	found.work, found.roster, found.missing = best, roster, false
	return found, nil
}

// pickWork accepts only an exact match on some name the candidate actually
// carries. Fuzzy ranking is what puts a sticker on the wrong game.
//
// It matches in two passes because stripping all punctuation makes distinct
// titles collide: "1/2 summer+" folds to the same string as "1/2 summer", and
// the fan disc -- an entry with an empty roster -- won on search order alone.
// The first pass keeps punctuation and only folds case and spacing, which
// separates those two; the loose pass exists for the titles where the
// transcription and catalog disagree on decoration ("＊" vs "*", "～" vs "~").
func pickWork(candidates []catalogclient.Work, romaji, chinese string) *catalogclient.Work {
	for _, fold := range []func(string) string{tight, norm} {
		wanted := map[string]bool{}
		for _, name := range []string{romaji, chinese} {
			if n := fold(name); n != "" {
				wanted[n] = true
			}
		}
		var hit *catalogclient.Work
		for i := range candidates {
			for _, name := range workNames(&candidates[i]) {
				if !wanted[fold(name)] {
					continue
				}
				// Among equally exact matches prefer one that has characters:
				// a re-release or fan disc that catalog knows only as a title
				// cannot answer the character half of this backfill.
				if hit == nil {
					hit = &candidates[i]
				}
				break
			}
		}
		if hit != nil {
			return hit
		}
	}
	return nil
}

func workNames(w *catalogclient.Work) []string {
	names := []string{w.DisplayName}
	if w.Latin != nil {
		names = append(names, *w.Latin)
	}
	for _, text := range w.Localized {
		names = append(names, text.Value)
	}
	return names
}

func findInRoster(roster []catalogclient.Character, id int64) *catalogclient.Character {
	for i := range roster {
		if mustID(roster[i].ID) == id {
			return &roster[i]
		}
	}
	return nil
}

func pickCharacter(roster []catalogclient.Character, romaji, chinese string) *catalogclient.Character {
	wanted := map[string]bool{}
	for _, name := range []string{romaji, chinese} {
		if n := norm(name); n != "" {
			wanted[n] = true
		}
	}
	if len(wanted) == 0 {
		return nil
	}
	for i := range roster {
		names := []string{roster[i].DisplayName}
		if roster[i].Latin != nil {
			names = append(names, *roster[i].Latin)
		}
		for _, text := range roster[i].Localized {
			names = append(names, text.Value)
		}
		for _, name := range names {
			if wanted[norm(name)] {
				return &roster[i]
			}
		}
	}
	return nil
}

// tight folds only case and whitespace, so punctuation still separates two
// otherwise identical titles.
func tight(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// fillMissing keeps the transcription for the languages catalog does not have.
// catalog is the identity and always wins where it has an entry, but it knows
// no Chinese name for 27 of these characters, and showing ティナ・フルール・
// レニスフィア to a Chinese reader when 缇娜 was typed out by hand throws away
// good data. The snapshot is a display cache; this is what it is for.
func fillMissing(canonical, transcription map[string]string) map[string]string {
	for locale, value := range transcription {
		if value = strings.TrimSpace(value); value == "" {
			continue
		}
		if canonical[locale] == "" {
			canonical[locale] = value
		}
	}
	return canonical
}

// norm folds the differences that are noise between a transcription and a
// catalog row: case, spacing, and the punctuation that varies between a
// title's fullwidth original and its romaji ("＊" vs "*", "～" vs "~", the
// space in "紬 文德斯"). What survives is letters and digits.
func norm(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (b *backfill) syncSearchText(packID string) error {
	return b.db.Exec(`
		UPDATE pack p SET search_text = btrim(
			coalesce((SELECT string_agg(v, ' ') FROM jsonb_each_text(p.title) t(k, v)), '') || ' ' ||
			coalesce((SELECT string_agg(v, ' ') FROM jsonb_each_text(p.description) d(k, v)), '') || ' ' ||
			coalesce((SELECT string_agg(v, ' ') FROM jsonb_each_text(p.catalog_work_name) w(k, v)), '') || ' ' ||
			coalesce((SELECT string_agg(DISTINCT v, ' ') FROM sticker s, jsonb_each_text(s.catalog_work_name) sw(k, v)
			          WHERE s.pack_id = p.id), '')
		) WHERE p.id = ?`, packID).Error
}

func (b *backfill) report() {
	fmt.Printf("\nresolved %d stickers to a game, %d of them to a character\n", b.linked, b.charLinked)
	if !b.apply {
		fmt.Println("(dry run -- nothing was written; pass -apply)")
	}
	printCounts("games catalog could not confirm", b.unknownWork)
	printCounts("characters missing from their game's roster", b.unknownChar)
}

func printCounts(title string, counts map[string]int) {
	if len(counts) == 0 {
		return
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if counts[keys[i]] != counts[keys[j]] {
			return counts[keys[i]] > counts[keys[j]]
		}
		return keys[i] < keys[j]
	})
	fmt.Printf("\n%s (%d):\n", title, len(keys))
	for _, key := range keys {
		fmt.Printf("  %3d  %s\n", counts[key], key)
	}
}

func label(romaji, chinese string) string {
	switch {
	case romaji != "" && chinese != "" && romaji != chinese:
		return romaji + " / " + chinese
	case romaji != "":
		return romaji
	default:
		return chinese
	}
}

func decode(raw datatypes.JSON) map[string]string {
	out := map[string]string{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out
}

func encode(in map[string]string) datatypes.JSON {
	raw, err := json.Marshal(in)
	if err != nil {
		return datatypes.JSON("{}")
	}
	return datatypes.JSON(raw)
}

func mustID(raw string) int64 {
	var id int64
	if _, err := fmt.Sscanf(raw, "%d", &id); err != nil {
		fmt.Fprintf(os.Stderr, "unparseable catalog id %q\n", raw)
	}
	return id
}

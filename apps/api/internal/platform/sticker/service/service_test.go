package service

import (
	"strings"
	"testing"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/internal/platform/sticker/model"
	"kun-galgame-sticker-api/pkg/perm"

	"gorm.io/datatypes"
)

func TestViewerVisibility(t *testing.T) {
	owner := &model.Pack{OwnerUID: 7, Status: model.PackDraft}

	cases := []struct {
		name    string
		viewer  Viewer
		canSee  bool
		canEdit bool
	}{
		{"anonymous", Viewer{}, false, false},
		{"stranger", Viewer{UID: 8}, false, false},
		{"owner", Viewer{UID: 7}, true, true},
		{"moderator", Viewer{UID: 8, Roles: []string{perm.RoleModerator}}, true, true},
		{"creator is not a moderator", Viewer{UID: 8, Roles: []string{perm.RoleCreator}}, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.viewer.canSeeUnpublished(owner); got != tc.canSee {
				t.Errorf("canSeeUnpublished = %v, want %v", got, tc.canSee)
			}
			if got := tc.viewer.canEdit(owner); got != tc.canEdit {
				t.Errorf("canEdit = %v, want %v", got, tc.canEdit)
			}
		})
	}
}

func TestViewerUIDZeroIsNotAnOwner(t *testing.T) {
	// owner_uid is never 0, but an anonymous viewer's UID is. Comparing them
	// without the > 0 guard would hand every anonymous visitor ownership of
	// any row whose owner failed to resolve.
	orphan := &model.Pack{OwnerUID: 0, Status: model.PackDraft}
	if (Viewer{}).canEdit(orphan) {
		t.Fatal("an anonymous viewer must not be able to edit an owner-less pack")
	}
}

func TestSanitizeML(t *testing.T) {
	long := strings.Repeat("鲲", MaxTextRunes+50)

	got := decodeML(sanitizeML(dto.MultilingualText{
		"zh-cn":   " 表情包 ",
		"ZH-TW":   "表情包",
		"klingon": "nuqneH",
		"en-us":   "",
		"ja-jp":   long,
	}))

	if got["zh-cn"] != "表情包" {
		t.Errorf("zh-cn = %q, want trimmed", got["zh-cn"])
	}
	if got["zh-tw"] != "表情包" {
		t.Errorf("uppercase locale key was not normalized: %v", got)
	}
	if _, ok := got["klingon"]; ok {
		t.Error("unknown locale key survived")
	}
	if _, ok := got["en-us"]; ok {
		t.Error("empty value survived")
	}
	// Truncation is by rune. Slicing by byte, as the old code did, cut this
	// multi-byte string mid-character and produced invalid UTF-8.
	if runes := []rune(got["ja-jp"]); len(runes) != MaxTextRunes {
		t.Errorf("truncated to %d runes, want %d", len(runes), MaxTextRunes)
	}
	if !strings.HasSuffix(got["ja-jp"], "鲲") {
		t.Error("truncation split a character")
	}
}

func TestSearchTextCoversEveryLanguage(t *testing.T) {
	title := sanitizeML(dto.MultilingualText{"zh-cn": "猫娘", "ja-jp": "ねこ"})
	description := sanitizeML(dto.MultilingualText{"en-us": "cat girls"})

	got := searchText(title, description)
	for _, want := range []string{"猫娘", "ねこ", "cat girls"} {
		if !strings.Contains(got, want) {
			t.Errorf("search text %q is missing %q", got, want)
		}
	}
}

func TestIsImageHash(t *testing.T) {
	valid := strings.Repeat("a1", 32)
	cases := map[string]bool{
		valid:                         true,
		strings.ToUpper(valid):        false,
		strings.Repeat("a", 63):       false,
		strings.Repeat("a", 65):       false,
		strings.Repeat("g", 64):       false,
		"":                            false,
		strings.Repeat("a", 63) + " ": false,
	}
	for hash, want := range cases {
		if got := isImageHash(hash); got != want {
			t.Errorf("isImageHash(%.10q...) = %v, want %v", hash, got, want)
		}
	}
}

func TestRatingDefaultsToSFW(t *testing.T) {
	nsfw := model.RatingNSFW
	odd := int16(42)
	if rating(nil) != model.RatingSFW {
		t.Error("missing rating must default to sfw")
	}
	if rating(&odd) != model.RatingSFW {
		t.Error("an unknown rating must fall back to sfw, never to nsfw")
	}
	if rating(&nsfw) != model.RatingNSFW {
		t.Error("nsfw was not preserved")
	}
}

func TestOrderAlwaysHasATiebreaker(t *testing.T) {
	for _, sort := range []string{dto.SortNew, dto.SortHot, dto.SortUpdated, "garbage"} {
		if !strings.HasSuffix(orderFor(sort), "id DESC") {
			t.Errorf("orderFor(%q) = %q, must end with the id tiebreaker", sort, orderFor(sort))
		}
	}
}

func jsonML(raw string) datatypes.JSON { return datatypes.JSON(raw) }

func contains(haystack, needle string) bool { return strings.Contains(haystack, needle) }

func count(haystack, needle string) int { return strings.Count(haystack, needle) }

package service

import (
	"encoding/json"
	"testing"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/internal/platform/sticker/model"
	"kun-galgame-sticker-api/internal/platform/sticker/repository"
)

func TestFaceRatingSpeaksCatalogVocabulary(t *testing.T) {
	if got := faceRating(model.RatingSFW); got != "all_ages" {
		t.Errorf("sfw = %q, want all_ages", got)
	}
	if got := faceRating(model.RatingNSFW); got != "r18" {
		t.Errorf("nsfw = %q, want r18", got)
	}
}

// The published contract must not grow fields by accident: the face DTO is
// built from the site DTO, and the site DTO carries things -- a draft status,
// an owner's uid on an unpublished row -- that the face has no business
// emitting. Marshalling one and reading its keys is the cheapest guard.
func TestFacePackEmitsNoStatus(t *testing.T) {
	raw, err := json.Marshal(facePack(dto.Pack{
		ID: "01a0", Status: model.PackDraft, ContentRating: model.RatingNSFW,
		Title: dto.MultilingualText{"zh-cn": "测试"},
	}))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, forbidden := range []string{"status", "is_official", "search_text", "owner_uid"} {
		if _, found := keys[forbidden]; found {
			t.Errorf("face pack leaks %q", forbidden)
		}
	}
	for _, required := range []string{"object", "id", "title", "content_rating", "official"} {
		if _, found := keys[required]; !found {
			t.Errorf("face pack is missing %q", required)
		}
	}
}

func TestFaceListPageArithmetic(t *testing.T) {
	cases := []struct {
		offset, limit, want int
	}{
		{0, 20, 1},
		{20, 20, 2},
		{100, 50, 3},
		{0, 0, 1},
	}
	for _, tc := range cases {
		if got := pageOf(indexParams(tc.offset, tc.limit)); got != tc.want {
			t.Errorf("offset %d limit %d = page %d, want %d", tc.offset, tc.limit, got, tc.want)
		}
	}
}

func indexParams(offset, limit int) repository.IndexParams {
	return repository.IndexParams{Offset: offset, Limit: limit}
}

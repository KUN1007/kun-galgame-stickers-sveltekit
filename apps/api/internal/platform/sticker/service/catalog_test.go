package service

import (
	"testing"

	"kun-galgame-sticker-api/pkg/catalogclient"
)

func TestCatalogNameMapsLocalesOntoTheSiteKeySpace(t *testing.T) {
	latin := "Suzuma Yukariko"
	got := catalogName("鈴間紫子", map[string]catalogclient.LocalizedText{
		"zh-Hans": {Value: "铃间紫子"},
		"zh-Hant": {Value: "鈴間紫子"},
		"ja":      {Value: "すずま ゆかりこ"},
		// catalog carries tags this site has no column for; they are dropped
		// rather than stored under a key no page renders.
		"ko": {Value: "스즈마"},
	}, &latin)

	want := map[string]string{
		"und":   "鈴間紫子",
		"zh-cn": "铃间紫子",
		"zh-tw": "鈴間紫子",
		"ja-jp": "すずま ゆかりこ",
		"en-us": "Suzuma Yukariko",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d keys, want %d: %v", len(got), len(want), got)
	}
	for key, value := range want {
		if got[key] != value {
			t.Errorf("%s = %q, want %q", key, got[key], value)
		}
	}
}

// A localized value that only exists as an empty string must not shadow the
// canonical display name -- an empty title renders as a blank card.
func TestCatalogNameSkipsEmptyValues(t *testing.T) {
	got := catalogName("Summer Pockets", map[string]catalogclient.LocalizedText{
		"zh-Hans": {Value: "   "},
	}, nil)
	if _, ok := got["zh-cn"]; ok {
		t.Errorf("blank zh-Hans was stored: %v", got)
	}
	if got["und"] != "Summer Pockets" {
		t.Errorf("und = %q", got["und"])
	}
}

// latin fills en-us only when catalog has no English localization of its own.
func TestCatalogNameLatinDoesNotOverrideEnglish(t *testing.T) {
	latin := "Karenai Sekai to Owaru Hana"
	got := catalogName("枯れない世界と終わる花", map[string]catalogclient.LocalizedText{
		"en": {Value: "The Unwithering World and the Ending Flower"},
	}, &latin)
	if got["en-us"] != "The Unwithering World and the Ending Flower" {
		t.Errorf("en-us = %q, want the localized English", got["en-us"])
	}
}

func TestParseCatalogIDRejectsNonPositive(t *testing.T) {
	for _, raw := range []string{"", "0", "-3", "abc", "12x"} {
		if got := parseCatalogID(raw); got != nil {
			t.Errorf("parseCatalogID(%q) = %d, want nil", raw, *got)
		}
	}
	if got := parseCatalogID(" 42 "); got == nil || *got != 42 {
		t.Errorf("parseCatalogID(\" 42 \") did not parse")
	}
}

// searchText is what the trigram index sits on: every language of every field,
// deduplicated so a pack that repeats its game name does not weight it twice.
func TestSearchTextFlattensAndDeduplicates(t *testing.T) {
	title := jsonML(`{"zh-cn":"水星的日常表情","en-us":"Mercury dailies"}`)
	description := jsonML(`{"zh-cn":"水星的日常表情"}`)
	work := jsonML(`{"und":"Summer Pockets","zh-cn":"夏日口袋"}`)

	got := searchText(title, description, work)
	for _, want := range []string{"水星的日常表情", "Mercury dailies", "Summer Pockets", "夏日口袋"} {
		if !contains(got, want) {
			t.Errorf("search text %q is missing %q", got, want)
		}
	}
	if count(got, "水星的日常表情") != 1 {
		t.Errorf("duplicate value repeated in %q", got)
	}
}

// catalog files character romanizations under lang=ja. Taking that at face
// value shows a Japanese reader "Kushima Kamome" instead of 久島鴎, which the
// response already carries as the canonical name.
func TestCatalogNameKeepsRomanizationOutOfTheJapaneseSlot(t *testing.T) {
	got := catalogName("久島鴎", map[string]catalogclient.LocalizedText{
		"ja":      {Value: "Kushima Kamome"},
		"zh-Hans": {Value: "久岛鸥"},
	}, nil)

	if got["ja-jp"] != "" {
		t.Errorf("ja-jp = %q, want the romanization to stay out of it", got["ja-jp"])
	}
	if got["en-us"] != "Kushima Kamome" {
		t.Errorf("en-us = %q, want the romanization", got["en-us"])
	}
	if got["und"] != "久島鴎" {
		t.Errorf("und = %q, want the canonical name", got["und"])
	}
	if got["zh-cn"] != "久岛鸥" {
		t.Errorf("zh-cn = %q", got["zh-cn"])
	}
}

// A genuine Japanese title under lang=ja still lands in the Japanese slot --
// the rule keys on script, not on the tag being suspicious.
func TestCatalogNameKeepsRealJapaneseTitles(t *testing.T) {
	got := catalogName("枯れない世界と終わる花", map[string]catalogclient.LocalizedText{
		"ja":      {Value: "枯れない世界と終わる花"},
		"zh-Hans": {Value: "永不枯萎的世界与终之花"},
	}, nil)
	if got["ja-jp"] != "枯れない世界と終わる花" {
		t.Errorf("ja-jp = %q", got["ja-jp"])
	}
}

// A Latin display name with a Latin ja value is left alone: there is no
// canonical Japanese form to prefer, so the slot keeps what catalog sent.
func TestCatalogNameLeavesLatinTitlesAlone(t *testing.T) {
	got := catalogName("Summer Pockets", map[string]catalogclient.LocalizedText{
		"ja": {Value: "Summer Pockets"},
	}, nil)
	if got["ja-jp"] != "Summer Pockets" {
		t.Errorf("ja-jp = %q, want the value kept", got["ja-jp"])
	}
}

// catalog rates images safe | suggestive | explicit but leaves nearly every row
// unassessed, so the gate must drop exactly the explicit ones and keep the rest
// -- requiring "safe" would blank almost every portrait on the site.
func TestImageURLDropsOnlyExplicit(t *testing.T) {
	safe, suggestive, explicit := "safe", "suggestive", "explicit"
	cases := []struct {
		name   string
		sexual *string
		want   string
	}{
		{"unassessed", nil, "https://cdn/x.webp"},
		{"safe", &safe, "https://cdn/x.webp"},
		{"suggestive", &suggestive, "https://cdn/x.webp"},
		{"explicit", &explicit, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := imageURL(&catalogclient.Image{URL: "https://cdn/x.webp", Sexual: tc.sexual})
			if got != tc.want {
				t.Errorf("imageURL = %q, want %q", got, tc.want)
			}
		})
	}
	if got := imageURL(nil); got != "" {
		t.Errorf("imageURL(nil) = %q", got)
	}
}

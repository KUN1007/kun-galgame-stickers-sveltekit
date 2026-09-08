package catalogclient

import "strings"

// siteLocales maps catalog's BCP-47 tags onto the key space this site's
// database has always used. A tag outside the map is dropped rather than
// stored under a key no page knows how to render.
var siteLocales = map[string]string{
	"zh-Hans": "zh-cn",
	"zh-Hant": "zh-tw",
	"zh":      "zh-cn",
	"ja":      "ja-jp",
	"en":      "en-us",
}

// LocalizedName folds a catalog entity's display name and localized variants
// into one multilingual value. display_name is the canonical name (usually
// Japanese) and is kept under "und" so a locale with no translation still
// renders something instead of an empty card.
func LocalizedName(displayName string, localized map[string]LocalizedText, latin *string) map[string]string {
	out := map[string]string{}
	displayName = strings.TrimSpace(displayName)
	if displayName != "" {
		out["und"] = displayName
	}
	for tag, text := range localized {
		key, ok := siteLocales[tag]
		if !ok {
			continue
		}
		value := strings.TrimSpace(text.Value)
		if value == "" {
			continue
		}
		// Character aliases file their romanization under lang=ja, so catalog
		// answers 久島鴎 with localized["ja"] = "Kushima Kamome". Rendering that
		// to a Japanese reader is worse than the canonical name they already
		// have, so a Latin-script value never takes the Japanese slot -- it
		// becomes the Latin form instead, which is what it is.
		if key == "ja-jp" && !HasJapaneseScript(value) && HasJapaneseScript(displayName) {
			if out["en-us"] == "" {
				out["en-us"] = value
			}
			continue
		}
		out[key] = value
	}
	if latin != nil {
		if value := strings.TrimSpace(*latin); value != "" && out["en-us"] == "" {
			out["en-us"] = value
		}
	}
	return out
}

// HasJapaneseScript reports whether a string contains kana or CJK ideographs.
// It is deliberately coarse: it only has to tell a romanization apart from a
// name written in the language it claims.
func HasJapaneseScript(value string) bool {
	for _, r := range value {
		switch {
		case r >= 0x3040 && r <= 0x30FF, // hiragana + katakana
			r >= 0x4E00 && r <= 0x9FFF, // CJK unified ideographs
			r >= 0x3400 && r <= 0x4DBF: // CJK extension A
			return true
		}
	}
	return false
}

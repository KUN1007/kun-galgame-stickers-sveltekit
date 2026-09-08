package dto

type MultilingualText map[string]string

type Author struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Tag struct {
	ID        string           `json:"id"`
	Slug      string           `json:"slug"`
	Name      MultilingualText `json:"name"`
	PackCount int              `json:"pack_count"`
}

// CatalogWork and CatalogCharacter are the snapshot this site stores of an
// infra catalog identity. ID is catalog's, and is what a link resolves against.
type CatalogWork struct {
	ID          int64            `json:"id"`
	Name        MultilingualText `json:"name"`
	CoverURL    string           `json:"cover_url"`
	ReleaseDate *string          `json:"release_date,omitempty"`
	Medium      string           `json:"medium,omitempty"`
}

type CatalogCharacter struct {
	ID         int64            `json:"id"`
	Name       MultilingualText `json:"name"`
	ImageURL   string           `json:"image_url"`
	RosterRole string           `json:"roster_role,omitempty"`
	Gender     *string          `json:"gender,omitempty"`
	Birthday   *string          `json:"birthday,omitempty"`
	BloodType  *string          `json:"blood_type,omitempty"`
	Traits     []CatalogTrait   `json:"traits,omitempty"`
	Aliases    []string         `json:"aliases,omitempty"`
}

type CatalogTrait struct {
	Name  MultilingualText `json:"name"`
	Group MultilingualText `json:"group"`
}

type Sticker struct {
	ID            string           `json:"id"`
	PackID        string           `json:"pack_id"`
	Position      int              `json:"position"`
	Width         int              `json:"width"`
	Height        int              `json:"height"`
	Game          MultilingualText `json:"game"`
	CharacterName MultilingualText `json:"character_name"`
	VndbID        *int             `json:"vndb_id,omitempty"`
	Note          string           `json:"note"`
	ImageURL      string           `json:"image_url"`
	ThumbURL      string           `json:"thumb_url"`

	CatalogWork      *CatalogWork      `json:"catalog_work,omitempty"`
	CatalogCharacter *CatalogCharacter `json:"catalog_character,omitempty"`
}

type Pack struct {
	ID            string           `json:"id"`
	Status        int16            `json:"status"`
	IsOfficial    bool             `json:"is_official"`
	ContentRating int16            `json:"content_rating"`
	Title         MultilingualText `json:"title"`
	Description   MultilingualText `json:"description"`
	CoverURL      string           `json:"cover_url"`
	CoverThumbURL string           `json:"cover_thumb_url"`
	StickerCount  int              `json:"sticker_count"`
	ViewCount     int64            `json:"view_count"`
	DownloadCount int64            `json:"download_count"`
	Author        Author           `json:"author"`
	Tags          []Tag            `json:"tags"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
	PublishedAt   *string          `json:"published_at,omitempty"`

	CatalogWork *CatalogWork `json:"catalog_work,omitempty"`
}

type PackDetail struct {
	Pack
	Stickers []Sticker `json:"stickers"`
	// Characters are the distinct catalog characters across this pack's
	// stickers, and Works the distinct games. A pack that declares one game is
	// the common case, but the seeded official packs span dozens.
	Characters []CatalogCharacter `json:"characters"`
	Works      []CatalogWork      `json:"works"`
}

// CharacterPage is the public /character/{id} route: catalog's profile plus
// every sticker on this site that is tagged with that character.
type CharacterPage struct {
	Character   CatalogCharacter   `json:"character"`
	Stickers    []CharacterSticker `json:"stickers"`
	Packs       map[string]Pack    `json:"packs"`
	Appearances []CatalogWork      `json:"appearances"`
	// Profile says whether the catalog half loaded. False means the page is
	// rendering from stored snapshots because catalog was unreachable.
	Profile bool `json:"profile"`
}

type CharacterSticker struct {
	ID       string `json:"id"`
	PackID   string `json:"pack_id"`
	ImageURL string `json:"image_url"`
	ThumbURL string `json:"thumb_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type PackListPage struct {
	Packs []Pack `json:"packs"`
	Total int64  `json:"total"`
}

type UploadResult struct {
	Hash     string `json:"hash"`
	ImageURL string `json:"image_url"`
	ThumbURL string `json:"thumb_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

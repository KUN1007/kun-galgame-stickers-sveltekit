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
	// ContentRating is catalog's all_ages | sensitive | r18. Unlike the images'
	// sexual field it is populated, so it is what the UI badges.
	ContentRating string `json:"content_rating,omitempty"`
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
	// Set on search hits: the game the character is from, and how many
	// stickers of them this site has, so a palette row can say where it leads.
	WorkName     MultilingualText `json:"work_name,omitempty"`
	StickerCount int              `json:"sticker_count,omitempty"`
}

// SearchResults is the command palette's payload: two short lanes, each
// already shaped the way its card renders.
type SearchResults struct {
	Packs      []Pack             `json:"packs"`
	Characters []CatalogCharacter `json:"characters"`
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

// Comment is one post in a pack's comment thread. The body arrives already
// cooked and sanitized by community; this site renders content_html and never
// re-processes it.
type Comment struct {
	ID          int64  `json:"id"`
	PostNumber  int    `json:"post_number"`
	ContentHTML string `json:"content_html"`
	ContentRaw  string `json:"content_raw"`
	CreatedAt   string `json:"created_at"`
	EditedAt    string `json:"edited_at,omitempty"`
	Author      Author `json:"author"`
	// Told to the client so the UI does not offer an action the API will
	// refuse; the API checks again regardless.
	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
	LikeCount int  `json:"like_count"`
	IsLiked   bool `json:"is_liked"`
	// ReplyTo is what the author answered; RootID is the top-level comment the
	// exchange hangs under, which community derives rather than the caller.
	ReplyTo     int64  `json:"reply_to,omitempty"`
	RootID      int64  `json:"root_id,omitempty"`
	ReplyToName string `json:"reply_to_name,omitempty"`
}

type CommentPage struct {
	ThreadID   int64     `json:"thread_id"`
	Comments   []Comment `json:"comments"`
	Total      int       `json:"total"`
	NextCursor string    `json:"next_cursor,omitempty"`
	// Enabled is false when the community service is not configured, which is
	// how a pack page knows to leave the section out rather than show an error.
	Enabled bool `json:"enabled"`
}

type CommentRequest struct {
	Body    string `json:"body"`
	ReplyTo int64  `json:"reply_to"`
}

type CommentFlagRequest struct {
	Reason int    `json:"reason"`
	Note   string `json:"note"`
}

type CommentLikeResult struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}

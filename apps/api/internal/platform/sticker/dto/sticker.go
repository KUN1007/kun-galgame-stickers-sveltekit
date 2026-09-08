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
}

type PackDetail struct {
	Pack
	Stickers []Sticker `json:"stickers"`
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

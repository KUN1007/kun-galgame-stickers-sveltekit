package dto

const (
	SortNew = "new"
	SortHot = "hot"
	// SortUpdated is only reachable from /me/packs: a draft has no
	// published_at, so the public sorts would bury every draft last.
	SortUpdated = "updated"
)

const (
	RatingFilterAll = "all"
	RatingFilterSFW = "sfw"
)

// ListQuery is already normalized when it reaches the service: the handler
// clamps page/limit and rejects unknown sort/rating values.
type ListQuery struct {
	Page         int
	Limit        int
	Sort         string
	Search       string
	Tag          string
	Rating       string
	OfficialOnly bool
	// LinkedOnly keeps only packs that declare a catalog work, so a reader
	// who came for galgame stickers can filter out everything else.
	LinkedOnly bool
	// CatalogWorkID lists the packs of one game.
	CatalogWorkID int64

	// OwnerUID scopes the list to one author. ViewerUID is who is asking --
	// only an author sees their own drafts.
	OwnerUID  int
	ViewerUID int
}

// A catalog id of 0 clears the link. JSON cannot tell an absent pointer from
// an explicit null once decoded, so a caller that means "unlink" says so with
// a value rather than by omitting the field.
type CreatePackRequest struct {
	Title         MultilingualText `json:"title"`
	Description   MultilingualText `json:"description"`
	ContentRating *int16           `json:"content_rating"`
	Tags          []string         `json:"tags"`
	CatalogWorkID *int64           `json:"catalog_work_id"`
}

type PatchPackRequest struct {
	Title          *MultilingualText `json:"title"`
	Description    *MultilingualText `json:"description"`
	ContentRating  *int16            `json:"content_rating"`
	CoverStickerID *string           `json:"cover_sticker_id"`
	Tags           *[]string         `json:"tags"`
	CatalogWorkID  *int64            `json:"catalog_work_id"`
}

type CreateStickerRequest struct {
	ImageHash          string           `json:"image_hash"`
	Width              int              `json:"width"`
	Height             int              `json:"height"`
	Game               MultilingualText `json:"game"`
	CharacterName      MultilingualText `json:"character_name"`
	VndbID             *int             `json:"vndb_id"`
	Note               *string          `json:"note"`
	CatalogWorkID      *int64           `json:"catalog_work_id"`
	CatalogCharacterID *int64           `json:"catalog_character_id"`
}

type PatchStickerRequest struct {
	Game               *MultilingualText `json:"game"`
	CharacterName      *MultilingualText `json:"character_name"`
	VndbID             *int              `json:"vndb_id"`
	Note               *string           `json:"note"`
	CatalogWorkID      *int64            `json:"catalog_work_id"`
	CatalogCharacterID *int64            `json:"catalog_character_id"`
}

type ReorderRequest struct {
	StickerIDs []string `json:"sticker_ids"`
}

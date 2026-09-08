-- Links packs and stickers to nextmoe-infra catalog identities.
--
-- catalog owns "who is who" (works, characters, external anchors); this site
-- keeps only the id plus a small display snapshot. The snapshot exists so list
-- pages never call upstream and the site still renders every name and portrait
-- while catalog is unreachable -- it is a cache, never the source of truth.
--
-- Both levels carry a work id on purpose. The seven official packs are mixed
-- collections (25, 32, 19... distinct games each), so a pack-level work cannot
-- describe them; per-sticker links stay accurate there, while a single-game
-- pack also declares its work once for the card and the discovery filter.
--
-- Names are JSONB in the same multilingual key space as pack.title
-- (zh-cn / zh-tw / ja-jp / en-us), so resolveMultilingual renders them like
-- every other translated field. Images are the catalog CDN URL verbatim.

ALTER TABLE pack
    ADD COLUMN catalog_work_id    BIGINT,
    ADD COLUMN catalog_work_name  JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN catalog_work_cover TEXT  NOT NULL DEFAULT '';

ALTER TABLE sticker
    ADD COLUMN catalog_work_id         BIGINT,
    ADD COLUMN catalog_work_name       JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN catalog_character_id    BIGINT,
    ADD COLUMN catalog_character_name  JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN catalog_character_image TEXT  NOT NULL DEFAULT '';

-- "Packs of this game" and the discovery "linked only" lane.
CREATE INDEX pack_catalog_work_idx ON pack (catalog_work_id, published_at DESC)
    WHERE catalog_work_id IS NOT NULL;
CREATE INDEX pack_linked_idx ON pack (published_at DESC)
    WHERE status = 1 AND catalog_work_id IS NOT NULL;

-- "Stickers of this character" powers the character page; the work index backs
-- the same lookup for packs whose stickers span several games.
CREATE INDEX sticker_catalog_character_idx ON sticker (catalog_character_id)
    WHERE catalog_character_id IS NOT NULL;
CREATE INDEX sticker_catalog_work_idx ON sticker (catalog_work_id)
    WHERE catalog_work_id IS NOT NULL;

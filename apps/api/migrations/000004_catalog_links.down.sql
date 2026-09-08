DROP INDEX IF EXISTS sticker_catalog_work_idx;
DROP INDEX IF EXISTS sticker_catalog_character_idx;
DROP INDEX IF EXISTS pack_linked_idx;
DROP INDEX IF EXISTS pack_catalog_work_idx;

ALTER TABLE sticker
    DROP COLUMN IF EXISTS catalog_character_image,
    DROP COLUMN IF EXISTS catalog_character_name,
    DROP COLUMN IF EXISTS catalog_character_id,
    DROP COLUMN IF EXISTS catalog_work_name,
    DROP COLUMN IF EXISTS catalog_work_id;

ALTER TABLE pack
    DROP COLUMN IF EXISTS catalog_work_cover,
    DROP COLUMN IF EXISTS catalog_work_name,
    DROP COLUMN IF EXISTS catalog_work_id;

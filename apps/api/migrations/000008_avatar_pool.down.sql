DROP INDEX IF EXISTS sticker_avatar_pool_slot_idx;
ALTER TABLE sticker DROP COLUMN IF EXISTS avatar_pool_slot;

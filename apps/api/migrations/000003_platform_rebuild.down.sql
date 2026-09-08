-- 000003 keeps the pre-rebuild rows in pack_legacy / sticker_legacy, so going
-- down restores the old schema with its data intact. Anything written through
-- the new API after the up ran is dropped here -- the legacy tables were never
-- written to again.

DROP TABLE IF EXISTS pack_tag;
DROP TABLE IF EXISTS tag;

-- pack and sticker reference each other (pack.cover_sticker_id / sticker.pack_id),
-- so neither can be dropped while the other stands. Break the cycle first.
ALTER TABLE IF EXISTS pack DROP CONSTRAINT IF EXISTS pack_cover_sticker_fk;
DROP TABLE IF EXISTS sticker;
DROP TABLE IF EXISTS pack;

ALTER TABLE pack_legacy    RENAME TO sticker_pack;
ALTER TABLE sticker_legacy RENAME TO sticker;

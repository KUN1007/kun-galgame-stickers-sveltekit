-- Trigram indexes for the quick-search palette.
--
-- pack.search_text already has one; the character lane had none because until
-- the catalog backfill there were no character names to search. The palette
-- queries both on every keystroke, so an unindexed ILIKE over every sticker is
-- the difference between a palette and a stutter.
--
-- The index is on the flattened JSONB text rather than a per-locale key: a
-- reader typing 鸣濑 and a reader typing Naruse must both find her, and the
-- names of both live in the same document.

CREATE INDEX sticker_catalog_character_name_idx
    ON sticker USING gin ((catalog_character_name::text) gin_trgm_ops);

CREATE INDEX sticker_catalog_work_name_idx
    ON sticker USING gin ((catalog_work_name::text) gin_trgm_ops);

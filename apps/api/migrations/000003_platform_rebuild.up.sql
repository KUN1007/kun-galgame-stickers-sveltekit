-- Rebuilds the domain around uuidv7 packs and stickers.
--
-- The old shape made sticker_pack.id serve as primary key, foreign key target
-- and public URL id all at once, and addressed a sticker as (sid, pid) where
-- pid came from max(pid)+1 -- two concurrent uploads raced onto the same pid
-- and the unique index turned that into a 500. Display order *was* pid, so a
-- pack could never be reordered and a deleted sticker left a permanent hole.
--
-- Existing rows: the 7 official packs (498 stickers, every image_hash already
-- backfilled into the infra image service) are copied over and marked
-- is_official. The legacy tables are kept, renamed, so this migration is
-- reversible; drop them in a follow-up once the new shape has run in
-- production for a release.
--
-- uuidv7() is a native PostgreSQL 18 function. Both dev and prod run
-- postgres:18-alpine.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

ALTER TABLE sticker      RENAME TO sticker_legacy;
ALTER TABLE sticker_pack RENAME TO pack_legacy;

DO $$
DECLARE
    missing INTEGER;
BEGIN
    SELECT count(*) INTO missing
    FROM sticker_legacy
    WHERE image_hash IS NULL OR btrim(image_hash) = '';

    IF missing > 0 THEN
        RAISE EXCEPTION
            'refusing to migrate: % legacy stickers have no image_hash; backfill them into the image service first',
            missing;
    END IF;
END $$;

CREATE TABLE pack (
    id               UUID        PRIMARY KEY DEFAULT uuidv7(),
    owner_uid        INTEGER     NOT NULL,
    status           SMALLINT    NOT NULL DEFAULT 0,
    is_official      BOOLEAN     NOT NULL DEFAULT FALSE,
    content_rating   SMALLINT    NOT NULL DEFAULT 0,
    title            JSONB       NOT NULL DEFAULT '{}',
    description      JSONB       NOT NULL DEFAULT '{}',
    cover_sticker_id UUID,
    sticker_count    INTEGER     NOT NULL DEFAULT 0,
    view_count       BIGINT      NOT NULL DEFAULT 0,
    download_count   BIGINT      NOT NULL DEFAULT 0,
    search_text      TEXT        NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at     TIMESTAMPTZ,
    legacy_id        INTEGER
);

CREATE TABLE sticker (
    id             UUID        PRIMARY KEY DEFAULT uuidv7(),
    pack_id        UUID        NOT NULL REFERENCES pack (id) ON DELETE CASCADE,
    position       INTEGER     NOT NULL,
    -- TEXT, not CHAR(64): the legacy column was CHAR and blank-padded every
    -- value, so every comparison needed a btrim first.
    image_hash     TEXT        NOT NULL,
    width          INTEGER     NOT NULL DEFAULT 0,
    height         INTEGER     NOT NULL DEFAULT 0,
    game           JSONB       NOT NULL DEFAULT '{}',
    character_name JSONB       NOT NULL DEFAULT '{}',
    vndb_id        INTEGER,
    note           TEXT        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE pack
    ADD CONSTRAINT pack_cover_sticker_fk
    FOREIGN KEY (cover_sticker_id) REFERENCES sticker (id) ON DELETE SET NULL;

INSERT INTO pack (
    owner_uid, status, is_official, title, description,
    sticker_count, search_text, created_at, updated_at, published_at, legacy_id
)
SELECT
    pl.owner_uid,
    pl.status,
    TRUE,
    pl.title,
    pl.description,
    (SELECT count(*) FROM sticker_legacy sl WHERE sl.sid = pl.id),
    btrim(
        coalesce((SELECT string_agg(v, ' ') FROM jsonb_each_text(pl.title)       AS t(k, v)), '') || ' ' ||
        coalesce((SELECT string_agg(v, ' ') FROM jsonb_each_text(pl.description) AS d(k, v)), '')
    ),
    pl.created,
    pl.updated,
    pl.published_at,
    pl.id
FROM pack_legacy pl
ORDER BY pl.id;

INSERT INTO sticker (
    pack_id, position, image_hash, game, character_name, vndb_id, note, created_at, updated_at
)
SELECT
    p.id,
    sl.pid,
    btrim(sl.image_hash),
    sl.game,
    sl.loli,
    NULLIF(sl.vndb, 0),
    sl.describe,
    sl.created,
    sl.updated
FROM sticker_legacy sl
JOIN pack p ON p.legacy_id = sl.sid
ORDER BY sl.sid, sl.pid;

-- position was copied straight from pid, so the old preview_pid identifies the
-- cover exactly. Matching on image_hash instead would be wrong: the image
-- service deduplicates, so one hash can belong to stickers in several packs.
UPDATE pack p
SET cover_sticker_id = s.id
FROM pack_legacy pl
JOIN sticker s ON s.position = pl.preview_pid
WHERE p.legacy_id = pl.id AND s.pack_id = p.id;

ALTER TABLE pack DROP COLUMN legacy_id;

CREATE INDEX pack_owner_idx     ON pack (owner_uid, created_at DESC);
CREATE INDEX pack_published_idx ON pack (published_at DESC) WHERE status = 1;
CREATE INDEX pack_hot_idx       ON pack (download_count DESC, view_count DESC) WHERE status = 1;
CREATE INDEX pack_official_idx  ON pack (published_at DESC) WHERE status = 1 AND is_official;
CREATE INDEX pack_search_idx    ON pack USING gin (search_text gin_trgm_ops);

CREATE UNIQUE INDEX sticker_pack_position_idx ON sticker (pack_id, position);
CREATE INDEX        sticker_hash_idx          ON sticker (image_hash);

CREATE TABLE tag (
    id         UUID    PRIMARY KEY DEFAULT uuidv7(),
    slug       TEXT    NOT NULL UNIQUE,
    name       JSONB   NOT NULL DEFAULT '{}',
    pack_count INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE pack_tag (
    pack_id UUID NOT NULL REFERENCES pack (id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES tag  (id) ON DELETE CASCADE,
    PRIMARY KEY (pack_id, tag_id)
);

CREATE INDEX pack_tag_tag_idx ON pack_tag (tag_id);

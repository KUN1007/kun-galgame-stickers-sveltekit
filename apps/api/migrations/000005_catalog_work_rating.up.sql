-- Carries catalog's content_rating into the display snapshot.
--
-- 89 of the 95 games these packs are made from are r18, so a game card that
-- cannot say so is not neutral -- a reader reasonably reads "no badge" as
-- all-ages. The rating comes from catalog with the name and cover and, like
-- them, is cached here so a list page renders the badge without an upstream
-- call.
--
-- Empty string, not NULL, for "not recorded": every other snapshot column on
-- these tables is NOT NULL DEFAULT, and one nullable text among them invites a
-- three-state check where the code only ever wants two.

ALTER TABLE pack
    ADD COLUMN catalog_work_rating TEXT NOT NULL DEFAULT '';

ALTER TABLE sticker
    ADD COLUMN catalog_work_rating TEXT NOT NULL DEFAULT '';

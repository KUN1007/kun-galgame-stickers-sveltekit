-- Local mirror of comment likes.
--
-- The community service owns whether a reaction is on or off -- its toggle is
-- the source of truth and it feeds the trust engine -- but its post projection
-- carries no reaction fields at all, so a consuming site cannot render "3
-- likes" or "you liked this" from a read. Every other site on this platform
-- solves it the same way: mirror the toggle's outcome locally and count here.
--
-- post_id is community's id, not this site's, which is why it is a plain
-- bigint with no foreign key: the row it points at lives in another database.

CREATE TABLE comment_like (
    post_id    BIGINT      NOT NULL,
    user_id    INTEGER     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, user_id)
);

-- "Which of these posts did I like" is the per-page query the pack detail runs.
CREATE INDEX comment_like_user_idx ON comment_like (user_id, post_id);

package repository

import (
	"kun-galgame-sticker-api/internal/platform/sticker/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentLikeRepo struct{ db *gorm.DB }

func NewCommentLikeRepo(db *gorm.DB) *CommentLikeRepo { return &CommentLikeRepo{db: db} }

// Ensure and Remove mirror community's toggle outcome. Both are idempotent:
// community has already decided the new state, and this side only has to end
// up agreeing with it.
func (r *CommentLikeRepo) Ensure(postID int64, userID int) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model.CommentLike{PostID: postID, UserID: userID}).Error
}

func (r *CommentLikeRepo) Remove(postID int64, userID int) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&model.CommentLike{}).Error
}

// Counts and LikedSet answer a whole page in one query each, so a thread of
// thirty comments costs two lookups rather than sixty.
func (r *CommentLikeRepo) Counts(postIDs []int64) (map[int64]int, error) {
	out := make(map[int64]int, len(postIDs))
	if len(postIDs) == 0 {
		return out, nil
	}
	var rows []struct {
		PostID int64
		N      int
	}
	err := r.db.Model(&model.CommentLike{}).
		Select("post_id, count(*) AS n").
		Where("post_id IN ?", postIDs).
		Group("post_id").Scan(&rows).Error
	for _, row := range rows {
		out[row.PostID] = row.N
	}
	return out, err
}

func (r *CommentLikeRepo) LikedSet(userID int, postIDs []int64) (map[int64]bool, error) {
	out := make(map[int64]bool, len(postIDs))
	if userID <= 0 || len(postIDs) == 0 {
		return out, nil
	}
	var ids []int64
	err := r.db.Model(&model.CommentLike{}).
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Pluck("post_id", &ids).Error
	for _, id := range ids {
		out[id] = true
	}
	return out, err
}

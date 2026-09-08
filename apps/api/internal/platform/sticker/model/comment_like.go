package model

import "time"

// CommentLike mirrors a community reaction so this site can render counts.
// post_id belongs to the community service, not to this database, which is why
// there is no foreign key to hang it on.
type CommentLike struct {
	PostID    int64     `gorm:"column:post_id;primaryKey"`
	UserID    int       `gorm:"column:user_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (CommentLike) TableName() string { return "comment_like" }

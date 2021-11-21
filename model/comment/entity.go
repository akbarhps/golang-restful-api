package comment

import "time"

type Comment struct {
	CommentID int64  `gorm:"column:comment_id;primaryKey,autoIncrement"`
	Content   string `gorm:"column:content"`
	PostID    string `gorm:"column:post_id"`
	UserID    string `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

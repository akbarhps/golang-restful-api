package like

import "time"

type Like struct {
	ID     int64  `gorm:"column:like_id;primaryKey,autoIncrement"`
	PostID string `gorm:"column:post_id"`
	UserID    string    `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

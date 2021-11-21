package post

import (
	"time"
)

type Post struct {
	ID     string `gorm:"column:post_id;primaryKey"`
	UserID string `gorm:"column:user_id;"`
	Caption   string    `gorm:"column:caption;"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
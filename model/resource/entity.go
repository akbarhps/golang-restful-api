package resource

import "time"

type Resource struct {
	ID          string `gorm:"column:resource_id"`
	IndexInPost int    `gorm:"column:index_in_post"`
	Path        string    `gorm:"column:path"`
	ShareURL    string    `gorm:"column:share_url"`
	PostID      string    `gorm:"column:post_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

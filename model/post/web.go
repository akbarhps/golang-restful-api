package post

import (
	"go-api/model/resource"
	"time"
)

type (
	CreateRequest struct {
		Caption   string              `json:"caption" form:"caption"`
		Resources []resource.Resource `json:"resources"`
		UserID    string              `json:"user_id"`
	}

	FindByPostIDRequest struct {
		PostID string `json:"post_id"`
		UserID string `json:"user_id"`
	}

	UpdateRequest struct {
		PostID  string `validate:"required" json:"post_id"`
		Caption string `validate:"required" json:"caption"`
		UserID  string `json:"user_id"`
	}

	DeleteRequest struct {
		PostID string `validate:"required" json:"post_id"`
		UserID string `validate:"required" json:"user_id"`
	}

	Response struct {
		PostID        string             `json:"post_id"`
		Thumbnail     *resource.Response `json:"thumbnail"`
		ResourceCount int64              `json:"resource_count"`
		LikesCount    int64              `json:"likes_count"`
		CommentsCount int64              `json:"comments_count"`
	}

	DetailResponse struct {
		PostID         string              `json:"post_id"`
		Caption        string              `json:"caption"`
		Resources      []resource.Response `json:"resources"`
		LikesCount     int64               `json:"likes_count"`
		ViewerHasLiked bool                `json:"viewer_has_liked"`
		CommentsCount  int64               `json:"comments_count"`
		CreatedAt      time.Time           `json:"created_at"`
		UpdatedAt      time.Time           `json:"updated_at"`
	}
)

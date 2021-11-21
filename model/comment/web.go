package comment

type (
	CreateRequest struct {
		PostID  string `validate:"required" json:"post_id"`
		UserID  string `validate:"required" json:"user_id"`
		Content string `validate:"required" json:"content"`
	}

	DeleteRequest struct {
		PostID string `validate:"required" json:"post_id"`
		UserID string `validate:"required" json:"user_id"`
	}
)

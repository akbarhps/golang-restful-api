package like

type (
	Request struct {
		PostID string `validate:"required" json:"post_id"`
		UserID string `validate:"required" json:"user_id"`
	}
)

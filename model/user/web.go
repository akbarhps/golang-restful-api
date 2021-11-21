package user

type (
	RegisterRequest struct {
		Email       string `validate:"required,email" form:"email" json:"email"`
		Username    string `validate:"required,alphanum,max=18" form:"username" json:"username"`
		DisplayName string `validate:"required" form:"display_name" json:"display_name"`
		Password    string `validate:"required,min=8" form:"password" json:"password"`
	}

	LoginRequest struct {
		Handler  string `form:"handler" json:"handler"`
		Password string `validate:"required" form:"password" json:"password"`
	}

	UpdateProfileRequest struct {
		UserID      string `validate:"required,uuid4" form:"user_id" json:"user_id"`
		Email       string `validate:"required,email" form:"email" json:"email"`
		Username    string `validate:"required,alphanum,max=18" form:"username" json:"username"`
		DisplayName string `validate:"required" form:"display_name" json:"display_name"`
		Biography   string `form:"biography" json:"biography"`
	}

	UpdatePasswordRequest struct {
		UserID      string `validate:"required,uuid4" form:"user_id" json:"user_id"`
		OldPassword string `validate:"required,min=8" form:"old_password" json:"old_password"`
		NewPassword string `validate:"required,min=8" form:"new_password" json:"new_password"`
	}

	AuthResponse struct {
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}

	SearchResponse struct {
		Username          string `json:"username"`
		DisplayName       string `json:"display_name"`
		ProfilePictureURL string `json:"profile_picture_url"`
	}

	Response struct {
		Username          string `json:"username"`
		DisplayName       string `json:"display_name"`
		Biography         string `json:"biography"`
		ExternalUrl       string `json:"external_url"`
		ProfilePictureURL string `json:"profile_picture_url"`
		IsVerified        bool   `json:"is_verified"`
		FollowedByViewer  bool   `json:"followed_by_viewer"`
		FollowerCount     int64  `json:"follower_count"`
		FollowingCount    int64  `json:"following_count"`
	}
)

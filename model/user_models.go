package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserRegisterRequest struct {
		FullName string `validate:"required,max=200,min=1" json:"full_name"`
		Username string `validate:"required,min=1,max=16" json:"username"`
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required,min=6" json:"password"`
	}

	UserLoginRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `validate:"required,min=6" json:"password"`
	}

	UserUpdateProfileRequest struct {
		Id       uuid.UUID `json:"id"`
		FullName string    `validate:"required,max=200,min=1" json:"full_name"`
		Username string    `validate:"required,min=1,max=16" json:"username"`
		Email    string    `validate:"required,email" json:"email"`
	}

	UserUpdatePasswordRequest struct {
		Id          uuid.UUID `json:"id"`
		Username    string    `validate:"required,min=1,max=16" json:"username"`
		Email       string    `validate:"required,email" json:"email"`
		OldPassword string    `validate:"required,min=6" json:"old_password"`
		NewPassword string    `validate:"required,min=6" json:"new_password"`
	}

	UserDeleteRequest struct {
		Id       uuid.UUID `json:"id"`
		Username string    `validate:"required,min=1,max=16" json:"username"`
		Email    string    `validate:"required,email" json:"email"`
		Password string    `validate:"required,min=6" json:"password"`
	}

	UserFindRequest struct {
		FullName string `json:"full_name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	UserResponse struct {
		Id        uuid.UUID `json:"id"`
		FullName  string    `json:"full_name"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}
)

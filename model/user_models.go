package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserRegister struct {
		FullName string `validate:"required,max=200,min=1" json:"full_name"`
		Username string `validate:"required,min=1,max=16" json:"username"`
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required,min=6" json:"password"`
	}

	UserLogin struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `validate:"required,min=6" json:"password"`
	}

	UserUpdateProfile struct {
		Id       uuid.UUID `json:"id"`
		FullName string    `validate:"required,max=200,min=1" json:"full_name"`
		Username string    `validate:"required,min=1,max=16" json:"username"`
		Email    string    `validate:"required,email" json:"email"`
	}

	UserUpdatePassword struct {
		Id          uuid.UUID `json:"id"`
		Username    string    `validate:"required,min=1,max=16" json:"username"`
		Email       string    `validate:"required,email" json:"email"`
		OldPassword string    `validate:"required,min=6" json:"old_password"`
		NewPassword string    `validate:"required,min=6" json:"new_password"`
	}

	UserDelete struct {
		Id       uuid.UUID `json:"id"`
		Username string    `validate:"required,min=1,max=16" json:"username"`
		Email    string    `validate:"required,email" json:"email"`
		Password string    `validate:"required,min=6" json:"password"`
	}

	UserFind struct {
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

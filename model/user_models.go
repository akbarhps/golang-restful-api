package model

import (
	"time"
)

type (
	UserRegister struct {
		Email       string `validate:"required,email" json:"email"`
		Username    string `validate:"required,min=1,max=16" json:"username"`
		DisplayName string `validate:"required,max=200,min=1" json:"display_name"`
		Password    string `validate:"required,min=6" json:"password"`
	}

	UserLogin struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `validate:"required,min=6" json:"password"`
	}

	UserUpdateProfile struct {
		Id          string `json:"id"`
		Email       string `validate:"required,email" json:"email"`
		Username    string `validate:"required,min=1,max=16" json:"username"`
		DisplayName string `validate:"required,max=200,min=1" json:"display_name"`
	}

	UserUpdatePassword struct {
		Id          string `json:"id"`
		Email       string `validate:"required,email" json:"email"`
		Username    string `validate:"required,min=1,max=16" json:"username"`
		OldPassword string `validate:"required,min=6" json:"old_password"`
		NewPassword string `validate:"required,min=6" json:"new_password"`
	}

	UserDelete struct {
		Id       string `json:"id"`
		Username string `validate:"required,min=1,max=16" json:"username"`
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required,min=6" json:"password"`
	}

	UserFind struct {
		Email       string `json:"email"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
	}

	UserResponse struct {
		Id          string    `json:"id"`
		Email       string    `json:"email"`
		Username    string    `json:"username"`
		DisplayName string    `json:"display_name"`
		CreatedAt   time.Time `json:"created_at"`
	}
)

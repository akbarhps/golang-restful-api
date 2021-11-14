package domain

import (
	"go-api/model"
	"time"
)

type User struct {
	Id          string
	DisplayName string
	Username    string
	Email       string
	Password    string
	CreatedAt   time.Time
}

func (u *User) ToResponse() *model.UserResponse {
	return &model.UserResponse{
		Id:          u.Id,
		DisplayName: u.DisplayName,
		Username:    u.Username,
		Email:       u.Email,
		CreatedAt:   u.CreatedAt,
	}
}

package domain

import (
	"go-api/model"
	"time"
)

type User struct {
	Id        string
	FullName  string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (u *User) ToResponse() *model.UserResponse {
	return &model.UserResponse{
		Id:        u.Id,
		FullName:  u.FullName,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

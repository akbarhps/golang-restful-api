package service

import (
	"context"
	"go-api/model"
)

type UserService interface {
	Register(ctx context.Context, req *model.UserRegister) *model.UserResponse
	Login(ctx context.Context, req *model.UserLogin) *model.UserResponse
	Find(ctx context.Context, req *model.UserFind) []model.UserResponse
	UpdateProfile(ctx context.Context, req *model.UserUpdateProfile) *model.UserResponse
	UpdatePassword(ctx context.Context, req *model.UserUpdatePassword)
	Delete(ctx context.Context, req *model.UserDelete)
}

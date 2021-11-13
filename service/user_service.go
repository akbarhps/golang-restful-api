package service

import (
	"context"
	"go-api/model"
)

type UserService interface {
	Register(ctx context.Context, req *model.UserRegisterRequest) *model.UserResponse
	Login(ctx context.Context, req *model.UserLoginRequest) *model.UserResponse
	Find(ctx context.Context, req *model.UserFindRequest) []model.UserResponse
	UpdateProfile(ctx context.Context, req *model.UserUpdateProfileRequest) *model.UserResponse
	UpdatePassword(ctx context.Context, req *model.UserUpdatePasswordRequest)
	Delete(ctx context.Context, req *model.UserDeleteRequest)
}

package service

import (
	"context"
	"go-api/model"
)

type UserService interface {
	Register(ctx context.Context, req *model.UserRegisterRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserResponse, error)
	Find(ctx context.Context, req *model.UserFindRequest) ([]model.UserResponse, error)
	UpdateProfile(ctx context.Context, req *model.UserUpdateProfileRequest) (*model.UserResponse, error)
	UpdatePassword(ctx context.Context, req *model.UserUpdatePasswordRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, req *model.UserDeleteRequest) error
}

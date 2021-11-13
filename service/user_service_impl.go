package service

import (
	"context"
	"go-api/domain"
	"go-api/exception"
	"go-api/helper"
	"go-api/model"
	"go-api/repository"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	Validate   *validator.Validate
	Repository repository.UserRepository
}

func NewUserService(validate *validator.Validate, repository repository.UserRepository) UserService {
	return &UserServiceImpl{Validate: validate, Repository: repository}
}

func (service *UserServiceImpl) Register(ctx context.Context, req *model.UserRegister) *model.UserResponse {
	err := service.Validate.Struct(req)
	helper.PanicIfError(err)

	users := service.Repository.Find(ctx, &domain.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if len(users) > 0 {
		panic(exception.RecordDuplicateError{Message: "Username or Email already taken"})
	}

	uid, err := uuid.NewUUID()
	helper.PanicIfError(err)

	encrypt, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	helper.PanicIfError(err)

	user := &domain.User{
		Id:        uid,
		FullName:  req.FullName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	}

	service.Repository.Create(ctx, user)
	return user.ToResponse()
}

func (service *UserServiceImpl) Login(ctx context.Context, req *model.UserLogin) *model.UserResponse {
	err := service.Validate.Struct(req)
	helper.PanicIfError(err)

	users := service.Repository.Find(ctx, &domain.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if len(users) == 0 {
		panic(exception.RecordNotFoundError{Message: "Record not found"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(req.Password))
	helper.PanicIfError(err, exception.InvalidCredentialError{Message: "Invalid password"})

	return users[0].ToResponse()
}

func (service *UserServiceImpl) Find(ctx context.Context, req *model.UserFind) []model.UserResponse {
	users := service.Repository.Find(ctx, &domain.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
	})

	var responseUsers []model.UserResponse
	for _, user := range users {
		responseUsers = append(responseUsers, *user.ToResponse())
	}

	return responseUsers
}

func (service *UserServiceImpl) UpdateProfile(ctx context.Context, req *model.UserUpdateProfile) *model.UserResponse {
	err := service.Validate.Struct(req)
	helper.PanicIfError(err)

	users := service.Repository.Find(ctx, &domain.User{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
	})
	if len(users) == 0 {
		panic(exception.RecordNotFoundError{Message: "Record not found"})
	}

	user := &users[0]
	user.FullName = req.FullName
	user.Username = req.Username
	user.Email = req.Email

	service.Repository.Update(ctx, user)
	return user.ToResponse()
}

func (service *UserServiceImpl) UpdatePassword(ctx context.Context, req *model.UserUpdatePassword) {
	err := service.Validate.Struct(req)
	helper.PanicIfError(err)

	users := service.Repository.Find(ctx, &domain.User{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
	})
	if len(users) == 0 {
		panic(exception.RecordNotFoundError{Message: "Record not found"})
	}

	user := &users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	helper.PanicIfError(err, exception.InvalidCredentialError{Message: "Invalid old password"})

	encrypt, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	helper.PanicIfError(err)

	user.Password = string(encrypt)
	service.Repository.Update(ctx, user)
}

func (service *UserServiceImpl) Delete(ctx context.Context, req *model.UserDelete) {
	err := service.Validate.Struct(req)
	helper.PanicIfError(err)

	users := service.Repository.Find(ctx, &domain.User{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
	})
	if len(users) == 0 {
		panic(exception.RecordNotFoundError{Message: "Record not found"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(req.Password))
	helper.PanicIfError(err, exception.InvalidCredentialError{Message: "Invalid password"})

	service.Repository.Delete(ctx, &users[0])
}

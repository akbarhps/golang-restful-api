package service

import (
	"context"
	"fmt"
	"go-api/domain"
	"go-api/exception"
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

func (service *UserServiceImpl) Register(ctx context.Context, req *model.UserRegisterRequest) (*model.UserResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	users, err := service.Repository.Find(ctx, &domain.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, exception.RecordDuplicateError{Message: "Username or Email already taken"}
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	encrypt, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Id:        uid,
		FullName:  req.FullName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	}

	err = service.Repository.Create(ctx, user)
	if err != nil {
		return nil, exception.RecordDuplicateError{Message: err.Error()}
	}

	return user.ToResponse(), nil
}

func (service *UserServiceImpl) Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	users, err := service.Repository.Find(ctx, &domain.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, exception.RecordNotFoundError{Message: "Record not found"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(req.Password))
	if err != nil {
		return nil, exception.WrongCredentialError{Message: "Incorrect password"}
	}

	return users[0].ToResponse(), nil
}

func (service *UserServiceImpl) Find(ctx context.Context, req *model.UserFindRequest) ([]model.UserResponse, error) {
	user := &domain.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
	}

	users, err := service.Repository.Find(ctx, user)
	if err != nil {
		return nil, err
	}

	var responseUsers []model.UserResponse
	for _, user := range users {
		responseUsers = append(responseUsers, *user.ToResponse())
	}

	return responseUsers, nil
}

func (service *UserServiceImpl) UpdateProfile(ctx context.Context, req *model.UserUpdateProfileRequest) (*model.UserResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	users, err := service.Repository.Find(ctx, &domain.User{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, exception.RecordNotFoundError{Message: "Record not found"}
	}

	user := &users[0]
	user.FullName = req.FullName
	user.Username = req.Username
	user.Email = req.Email

	err = service.Repository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (service *UserServiceImpl) UpdatePassword(ctx context.Context, req *model.UserUpdatePasswordRequest) (*model.UserResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	users, err := service.Repository.Find(ctx, &domain.User{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		fmt.Println("Error di find")
		return nil, err
	}
	if len(users) == 0 {
		return nil, exception.RecordNotFoundError{Message: "Record not found"}
	}

	user := &users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return nil, exception.WrongCredentialError{Message: "Old password didn't match"}
	}

	encrypt, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error di generate")
		return nil, err
	}

	user.Password = string(encrypt)

	err = service.Repository.Update(ctx, user)
	if err != nil {
		fmt.Println("error di update")
		return nil, err
	}

	return user.ToResponse(), nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, req *model.UserDeleteRequest) error {
	err := service.Validate.Struct(req)
	if err != nil {
		return err
	}

	users, err := service.Repository.Find(ctx, &domain.User{
		Id:       req.Id,
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return exception.RecordNotFoundError{Message: "Record not found"}
	}

	user := &users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return exception.WrongCredentialError{Message: "Incorrect password"}
	}

	err = service.Repository.Delete(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"context"
	"database/sql"
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
	DB         *sql.DB
	Validate   *validator.Validate
	Repository repository.UserRepository
}

func NewUserService(DB *sql.DB, validate *validator.Validate, repository repository.UserRepository) UserService {
	return &UserServiceImpl{DB: DB, Validate: validate, Repository: repository}
}

func (s *UserServiceImpl) Register(ctx context.Context, req *model.UserRegister) *model.UserResponse {
	err := s.Validate.Struct(req)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.TXCommitOrRollback(tx)

	users := s.Repository.Find(ctx, tx, &domain.User{
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
		Id:        uid.String(),
		FullName:  req.FullName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(encrypt),
		CreatedAt: time.Now(),
	}

	s.Repository.Create(ctx, tx, user)
	return user.ToResponse()
}

func (s *UserServiceImpl) Login(ctx context.Context, req *model.UserLogin) *model.UserResponse {
	err := s.Validate.Struct(req)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.TXCommitOrRollback(tx)

	users := s.Repository.Find(ctx, tx, &domain.User{
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

func (s *UserServiceImpl) Find(ctx context.Context, req *model.UserFind) []model.UserResponse {
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.TXCommitOrRollback(tx)

	users := s.Repository.Find(ctx, tx, &domain.User{
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

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, req *model.UserUpdateProfile) *model.UserResponse {
	err := s.Validate.Struct(req)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.TXCommitOrRollback(tx)

	users := s.Repository.Find(ctx, tx, &domain.User{
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

	s.Repository.Update(ctx, tx, user)
	return user.ToResponse()
}

func (s *UserServiceImpl) UpdatePassword(ctx context.Context, req *model.UserUpdatePassword) {
	err := s.Validate.Struct(req)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.TXCommitOrRollback(tx)

	users := s.Repository.Find(ctx, tx, &domain.User{
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
	s.Repository.Update(ctx, tx, user)
}

func (s *UserServiceImpl) Delete(ctx context.Context, req *model.UserDelete) {
	err := s.Validate.Struct(req)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.TXCommitOrRollback(tx)

	users := s.Repository.Find(ctx, tx, &domain.User{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
	})
	if len(users) == 0 {
		panic(exception.RecordNotFoundError{Message: "Record not found"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(req.Password))
	helper.PanicIfError(err, exception.InvalidCredentialError{Message: "Invalid password"})

	s.Repository.Delete(ctx, tx, &users[0])
}

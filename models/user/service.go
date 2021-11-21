package user

import (
	"context"
	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"go-api/app"
	"go-api/exception"
	"go-api/helper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) *AuthResponse
	Login(ctx context.Context, req *LoginRequest) *AuthResponse

	UpdateProfile(ctx context.Context, req *UpdateProfileRequest)
	UpdatePassword(ctx context.Context, req *UpdatePasswordRequest)

	FindByUsername(ctx context.Context, username string) *Response
	SearchLike(ctx context.Context, keyword string) []*SearchResponse
}

type serviceImpl struct {
	validate       *validator.Validate
	userRepository Repository
}

func NewService(validate *validator.Validate, userRepository Repository) Service {
	return &serviceImpl{validate: validate, userRepository: userRepository}
}

func (s *serviceImpl) Register(ctx context.Context, req *RegisterRequest) *AuthResponse {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	var mErr exception.Errors

	fUser := s.userRepository.FindByUsername(tx, req.Username)
	if fUser.ID != "" {
		mErr.Errors = append(mErr.Errors, exception.FieldError{
			Field:   "username",
			Message: "username already taken",
		})
	}

	fUser = s.userRepository.FindByEmail(tx, req.Email)
	if fUser.ID != "" {
		mErr.Errors = append(mErr.Errors, exception.FieldError{
			Field:   "email",
			Message: "email already taken",
		})
	}

	if len(mErr.Errors) > 0 {
		panic(mErr)
	}

	encrypt, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	eUser := &User{
		ID:          uuid.NewV4().String(),
		Email:       req.Email,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Password:    string(encrypt),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	token, err := helper.GenerateJWT(eUser.ID, eUser.Username)
	if err != nil {
		panic(err)
	}

	s.userRepository.Create(tx, eUser)
	return &AuthResponse{
		UserID: eUser.ID,
		Token:  token,
	}
}

func (s *serviceImpl) Login(ctx context.Context, req *LoginRequest) *AuthResponse {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	user := s.userRepository.FindByEmailOrUsername(tx, req.Handler)
	if user.ID == "" {
		panic(exception.NotFoundError{
			Message: "username or email does not match any record",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		panic(exception.WrongPasswordError{Message: "password doest not match"})
	}

	token, err := helper.GenerateJWT(user.ID, user.Username)
	if err != nil {
		panic(err)
	}

	return &AuthResponse{
		UserID: user.ID,
		Token:  token,
	}
}

func (s *serviceImpl) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	user := s.userRepository.FindById(tx, req.UserID)
	if user.ID == "" {
		panic(exception.NotFoundError{
			Message: "user not found",
		})
	}

	var mErr exception.Errors
	user.Username = req.Username

	fUser := s.userRepository.FindByUsername(tx, req.Username)
	if fUser.ID != "" {
		mErr.Errors = append(mErr.Errors, exception.FieldError{
			Field:   "username",
			Message: "username already taken",
		})
	}

	user.Email = req.Email
	fUser = s.userRepository.FindByEmail(tx, req.Email)
	if fUser.ID != "" {
		mErr.Errors = append(mErr.Errors, exception.FieldError{
			Field:   "email",
			Message: "email already taken",
		})
	}

	if len(mErr.Errors) > 0 {
		panic(mErr)
	}

	user.DisplayName = req.DisplayName
	user.Biography = req.Biography
	s.userRepository.Update(tx, user)
}

func (s *serviceImpl) UpdatePassword(ctx context.Context, req *UpdatePasswordRequest) {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	user := s.userRepository.FindById(tx, req.UserID)
	if user.ID == "" {
		panic(exception.NotFoundError{
			Message: "user not found",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		panic(exception.WrongPasswordError{Message: "password doest not match"})
	}

	encrypt, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user.Password = string(encrypt)
	s.userRepository.Update(tx, user)
}

func (s *serviceImpl) FindByUsername(ctx context.Context, username string) *Response {
	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	user := s.userRepository.FindByUsername(tx, username)
	if user.ID == "" {
		panic(exception.NotFoundError{
			Message: "user not found",
		})
	}

	//TODO: Look for user followers and following
	//TODO: check if viewer following current user
	return &Response{
		Username:          user.Username,
		DisplayName:       user.DisplayName,
		Biography:         user.Biography,
		ExternalUrl:       user.ExternalUrl,
		ProfilePictureURL: user.Username,
		IsVerified:        user.IsVerified,
		FollowedByViewer:  false,
		FollowerCount:     0,
		FollowingCount:    0,
	}
}

func (s *serviceImpl) SearchLike(ctx context.Context, keyword string) []*SearchResponse {
	var sResponse []*SearchResponse
	tx := app.GetDB().WithContext(ctx).Begin()

	users := s.userRepository.FindLike(tx, keyword)
	for _, user := range users {
		sResponse = append(sResponse, &SearchResponse{
			Username:          user.Username,
			DisplayName:       user.DisplayName,
			ProfilePictureURL: user.Username,
		})
	}
	return sResponse
}

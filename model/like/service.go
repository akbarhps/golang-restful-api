package like

import (
	"context"
	"github.com/go-playground/validator"
	"go-api/app"
	"go-api/exception"
	"go-api/helper"
	"time"
)

type Service interface {
	Create(ctx context.Context, req *Request)
	Delete(ctx context.Context, req *Request)
}

type serviceImpl struct {
	validate *validator.Validate
	likeRepo Repository
}

func NewService(validate *validator.Validate, likeRepo Repository) Service {
	return &serviceImpl{validate: validate, likeRepo: likeRepo}
}

func (s *serviceImpl) Create(ctx context.Context, req *Request) {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	like := s.likeRepo.FindByPostIDAndUserID(tx, req.PostID, req.UserID)
	if like.LikeID != 0 {
		panic(exception.DuplicateError{Message:"can't like post more than once"})
	}

	s.likeRepo.Create(tx, &Like{
		PostID:    req.PostID,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
	})
}

func (s *serviceImpl) Delete(ctx context.Context, req *Request) {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	like := s.likeRepo.FindByPostIDAndUserID(tx, req.PostID, req.UserID)
	if like.LikeID == 0 {
		panic(exception.NotFoundError{Message: "like not found"})
	}

	s.likeRepo.Delete(tx, like.LikeID)
}

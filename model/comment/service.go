package comment

import (
	"context"
	"github.com/go-playground/validator"
	"go-api/app"
	"go-api/exception"
	"go-api/helper"
	"time"
)

type Service interface {
	Create(ctx context.Context, req *CreateRequest)
	Delete(ctx context.Context, req *DeleteRequest)
}

type serviceImpl struct {
	validate    *validator.Validate
	commentRepo Repository
}

func NewService(validate *validator.Validate, commentRepo Repository) Service {
	return &serviceImpl{validate: validate, commentRepo: commentRepo}
}

func (s *serviceImpl) Create(ctx context.Context, req *CreateRequest) {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	s.commentRepo.Create(tx, &Comment{
		Content:   req.Content,
		PostID:    req.PostID,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
	})
}

func (s *serviceImpl) Delete(ctx context.Context, req *DeleteRequest) {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	comment := s.commentRepo.FindByPostIDAndUserID(tx, req.PostID, req.UserID)
	if comment.CommentID == 0 {
		panic(exception.NotFoundError{Message: "comment not found"})
	}

	s.commentRepo.Delete(tx, comment.CommentID)
}

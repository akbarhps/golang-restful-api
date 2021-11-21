package post

import (
	"context"
	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"go-api/app"
	"go-api/exception"
	"go-api/helper"
	"go-api/model/comment"
	"go-api/model/like"
	"go-api/model/resource"
	"time"
)

type Service interface {
	Create(ctx context.Context, req *CreateRequest) *DetailResponse
	Update(ctx context.Context, req *UpdateRequest)
	Delete(ctx context.Context, req *DeleteRequest)
	FindByPostID(ctx context.Context, postID, viewerID string) *DetailResponse
	FindByUserID(ctx context.Context, userID string) []*Response
}

type serviceImpl struct {
	validate *validator.Validate
	postRepository     Repository
	resourceRepository resource.Repository
	likeRepository     like.Repository
	commentRepository  comment.Repository
}

func NewService(validate *validator.Validate, postRepository Repository, resourceRepository resource.Repository, likeRepository like.Repository, commentRepository comment.Repository) Service {
	return &serviceImpl{validate: validate, postRepository: postRepository, resourceRepository: resourceRepository, likeRepository: likeRepository, commentRepository: commentRepository}
}

func (s *serviceImpl) Create(ctx context.Context, req *CreateRequest) *DetailResponse {
	err := s.validate.Struct(req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	post := &Post{
		PostID:    uuid.NewV4().String(),
		Caption:   req.Caption,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.postRepository.Create(tx, post)

	var resourcesResp []resource.Response
	for _, r := range req.Resources {
		r.ResourceID = uuid.NewV4().String()
		r.CreatedAt = time.Now()
		r.PostID = post.PostID
		s.resourceRepository.Create(tx, &r)

		resourcesResp = append(resourcesResp, resource.Response{
			ShareURL: r.ShareURL,
		})
	}

	return &DetailResponse{
		PostID:    post.PostID,
		Caption:   post.Caption,
		Resources: resourcesResp,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func (s *serviceImpl) Update(ctx context.Context, req *UpdateRequest) {
	err := s.validate.Struct(&req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	fPost := s.postRepository.FindByPostID(tx, req.PostID)
	if fPost == nil {
		panic(exception.NotFoundError{Message: "post not found"})
	}

	if fPost.UserID != req.UserID {
		panic(exception.NoAccessError{Message: "can't update other person post"})
	}

	s.postRepository.Update(tx, &Post{
		PostID:    fPost.PostID,
		Caption:   req.Caption,
		UserID:    fPost.UserID,
		UpdatedAt: time.Now(),
	})
}

func (s *serviceImpl) Delete(ctx context.Context, req *DeleteRequest) {
	err := s.validate.Struct(&req)
	if err != nil {
		panic(err)
	}

	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	fPost := s.postRepository.FindByPostID(tx, req.PostID)
	if fPost == nil {
		panic(exception.NotFoundError{Message: "post not found"})
	}

	if fPost.UserID != req.UserID {
		panic(exception.NoAccessError{Message: "can't delete other person post"})
	}

	s.postRepository.Delete(tx, fPost.PostID)
}

func (s *serviceImpl) FindByPostID(ctx context.Context, postID, viewerID string) *DetailResponse {
	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	post := s.postRepository.FindByPostID(tx, postID)
	if post.PostID == "" {
		panic(exception.NotFoundError{Message: "post not found"})
	}

	var resResponse []resource.Response
	res := s.resourceRepository.FindByPostID(tx, postID)
	for _, r := range res {
		resResponse = append(resResponse, resource.Response{
			ShareURL: r.ShareURL,
		})
	}

	likesCount, hasViewerLiked := s.likeRepository.CountByPostID(tx, postID, viewerID)
	commentsCount := s.commentRepository.CountByPostID(tx, postID)
	return &DetailResponse{
		PostID:         post.PostID,
		Caption:        post.Caption,
		Resources:      resResponse,
		LikesCount:     likesCount,
		ViewerHasLiked: hasViewerLiked,
		CommentsCount:  commentsCount,
		CreatedAt:      post.CreatedAt,
		UpdatedAt:      post.UpdatedAt,
	}
}

func (s *serviceImpl) FindByUserID(ctx context.Context, userID string) []*Response {
	tx := app.GetDB().WithContext(ctx).Begin()
	defer helper.TXCommitOrRollback(tx)

	var response []*Response
	posts := s.postRepository.FindByUserID(tx, userID)
	for _, p := range posts {
		res, resCount := s.resourceRepository.FindFirstByPostID(tx, p.PostID)
		resourceResp := &resource.Response{ShareURL: res.ShareURL}

		likesCount, _ := s.likeRepository.CountByPostID(tx, p.PostID, userID)
		commentsCount := s.commentRepository.CountByPostID(tx, p.PostID)

		response = append(response, &Response{
			PostID:        p.PostID,
			Thumbnail:     resourceResp,
			ResourceCount: resCount,
			LikesCount:    likesCount,
			CommentsCount: commentsCount,
		})
	}
	return response
}

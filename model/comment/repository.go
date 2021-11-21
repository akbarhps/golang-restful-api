package comment

import (
	"go-api/exception"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, comment *Comment)
	Delete(tx *gorm.DB, commentID int64)
	CountByPostID(tx *gorm.DB, postID string) int64
	FindByPostID(tx *gorm.DB, postID string) []Comment
	FindByPostIDAndUserID(tx *gorm.DB, postID, userID string) *Comment
}

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (*repositoryImpl) Create(tx *gorm.DB, comment *Comment) {
	err := tx.Create(&comment).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Delete(tx *gorm.DB, commentID int64) {
	err := tx.Where("comment_id = ?", commentID).Delete(&Comment{}).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) CountByPostID(tx *gorm.DB, postID string) int64 {
	var commentsCount int64
	err := tx.Model(&Comment{}).
		Select("count(comment_id) as comments_count").
		Where("post_id = ?", postID).
		Find(&commentsCount).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return commentsCount
}

func (*repositoryImpl) FindByPostID(tx *gorm.DB, postID string) []Comment {
	var comments []Comment
	err := tx.Where("post_id = ?", postID).
		Order("created_at asc").
		Find(&comments).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return comments
}

func (*repositoryImpl) FindByPostIDAndUserID(tx *gorm.DB, postID, userID string) *Comment {
	var comment Comment
	err := tx.Where("post_id = ? and user_id = ?", postID, userID).
		Find(&comment).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return &comment
}

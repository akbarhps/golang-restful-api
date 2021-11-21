package post

import (
	"go-api/exception"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, post *Post)
	Update(tx *gorm.DB, post *Post)
	Delete(tx *gorm.DB, postID string)
	FindByPostID(tx *gorm.DB, postID string) *Post
	FindByUserID(tx *gorm.DB, userID string) []*Post
}

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (*repositoryImpl) Create(tx *gorm.DB, post *Post) {
	err := tx.Create(&post).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Update(tx *gorm.DB, post *Post) {
	err := tx.Model(&Post{}).
		Where("post_id = ? AND user_id = ?", post.PostID, post.UserID).
		Updates(&Post{
			Caption:   post.Caption,
			UpdatedAt: post.UpdatedAt,
		}).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Delete(tx *gorm.DB, postID string) {
	err := tx.Where("post_id = ?", postID).Delete(&Post{}).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) FindByPostID(tx *gorm.DB, postID string) *Post {
	var post *Post
	err := tx.Where("post_id = ?", postID).
		Limit(1).
		Find(&post).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return post
}

func (*repositoryImpl) FindByUserID(tx *gorm.DB, userID string) []*Post {
	var posts []*Post
	err := tx.Where("user_id = ?", userID).Find(&posts).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return posts
}

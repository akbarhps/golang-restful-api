package like

import (
	"go-api/exception"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, like *Like)
	Delete(tx *gorm.DB, likeID int64)
	CountByPostID(tx *gorm.DB, postID, userID string) (int64, bool)
	FindByPostID(tx *gorm.DB, postID string) []*Like
	FindByPostIDAndUserID(tx *gorm.DB, postID, userID string) *Like
}

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (*repositoryImpl) Create(tx *gorm.DB, like *Like) {
	err := tx.Create(&like).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Delete(tx *gorm.DB, likeID int64) {
	err := tx.Where("like_id = ?", likeID).Delete(&Like{}).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) CountByPostID(tx *gorm.DB, postID, userID string) (int64, bool) {
	var likesCount int64
	var viewerHasLiked bool
	err := tx.Model(&Like{}).
		Select("count(like_id) as likes_count").
		Select("if(strcmp(user_id, ?) = 0, 1, 0) as viewer_has_liked", userID).
		Where("post_id = ?", postID).
		Find(&likesCount, &viewerHasLiked).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return likesCount, viewerHasLiked
}

func (*repositoryImpl) FindByPostID(tx *gorm.DB, postID string) []*Like {
	var likes []*Like
	err := tx.Model(&Like{}).Where("post_id = ?", postID).Find(&likes).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return likes
}

func (*repositoryImpl) FindByPostIDAndUserID(tx *gorm.DB, postID, userID string) *Like {
	var like Like
	err := tx.Where("post_id = ? and user_id = ?", postID, userID).Find(&like).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return &like
}

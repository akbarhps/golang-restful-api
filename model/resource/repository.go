package resource

import (
	"go-api/exception"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, resource *Resource)
	Delete(tx *gorm.DB, resource *Resource)
	FindByResourceID(tx *gorm.DB, resourceID string) *Resource
	FindByPostID(tx *gorm.DB, postID string) []*Resource
	FindFirstByPostID(tx *gorm.DB, postID string) (*Resource, int64)
}

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (*repositoryImpl) Create(tx *gorm.DB, resource *Resource) {
	err := tx.Create(&resource).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Delete(tx *gorm.DB, resource *Resource) {
	err := tx.Delete(&resource).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) FindByResourceID(tx *gorm.DB, resourceID string) *Resource {
	var resource Resource
	err := tx.
		Where("resource_id = ?", resourceID).
		First(&resource).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return &resource
}

func (*repositoryImpl) FindByPostID(tx *gorm.DB, postID string) []*Resource {
	var resources []*Resource
	err := tx.
		Where("post_id = ?", postID).
		Order("created_at desc").
		Find(&resources).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return resources
}

func (*repositoryImpl) FindFirstByPostID(tx *gorm.DB, postID string) (*Resource, int64) {
	var resource Resource
	var resourcesCount int64
	err := tx.
		Where("post_id = ?", postID).
		Order("index_in_post asc").
		First(&resource).Count(&resourcesCount).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return &resource, resourcesCount
}

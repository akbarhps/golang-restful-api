package user

import (
	"go-api/exception"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(tx *gorm.DB, user *Entity)
	Update(tx *gorm.DB, user *Entity)
	Delete(tx *gorm.DB, user *Entity)
	FindById(tx *gorm.DB, id string) *Entity
	FindLike(tx *gorm.DB, keyword string) []*Entity
	FindByEmail(tx *gorm.DB, email string) *Entity
	FindByUsername(tx *gorm.DB, username string) *Entity
	FindByEmailOrUsername(tx *gorm.DB, handler string) *Entity
}

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (*repositoryImpl) Create(tx *gorm.DB, user *Entity) {
	err := tx.Create(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Update(tx *gorm.DB, user *Entity) {
	err := tx.Model(&Entity{}).
		Where("user_id = ?", user.ID).
		Updates(&Entity{
			DisplayName: user.DisplayName,
			Username:    user.Username,
			Email:       user.Email,
			Password:    user.Password,
			UpdatedAt:   time.Now(),
		}).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Delete(tx *gorm.DB, user *Entity) {
	err := tx.Where(&user).Delete(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) FindById(tx *gorm.DB, id string) *Entity {
	var user *Entity
	err := tx.Where("user_id = ?", id).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

func (*repositoryImpl) FindLike(tx *gorm.DB, keyword string) []*Entity {
	var users []*Entity
	query := "username LIKE ? OR display_name LIKE ?"
	key := "%" + keyword + "%"
	err := tx.Where(query, key, key).Find(&users).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return users
}

func (*repositoryImpl) FindByEmail(tx *gorm.DB, email string) *Entity {
	var user *Entity
	err := tx.Where("email = ?", email).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

func (*repositoryImpl) FindByUsername(tx *gorm.DB, username string) *Entity {
	var user *Entity
	err := tx.Where("username = ?", username).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

func (*repositoryImpl) FindByEmailOrUsername(tx *gorm.DB, handler string) *Entity {
	var user *Entity
	err := tx.Where("email = ? OR username = ?", handler, handler).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

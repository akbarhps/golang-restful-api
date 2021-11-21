package user

import (
	"go-api/exception"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(tx *gorm.DB, user *User)
	Update(tx *gorm.DB, user *User)
	Delete(tx *gorm.DB, user *User)
	FindById(tx *gorm.DB, id string) *User
	FindLike(tx *gorm.DB, keyword string) []*User
	FindByEmail(tx *gorm.DB, email string) *User
	FindByUsername(tx *gorm.DB, username string) *User
	FindByEmailOrUsername(tx *gorm.DB, handler string) *User
}

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (*repositoryImpl) Create(tx *gorm.DB, user *User) {
	err := tx.Create(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) Update(tx *gorm.DB, user *User) {
	err := tx.Model(&User{}).
		Where("user_id = ?", user.ID).
		Updates(&User{
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

func (*repositoryImpl) Delete(tx *gorm.DB, user *User) {
	err := tx.Where(&user).Delete(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
}

func (*repositoryImpl) FindById(tx *gorm.DB, id string) *User {
	var user *User
	err := tx.Where("user_id = ?", id).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

func (*repositoryImpl) FindLike(tx *gorm.DB, keyword string) []*User {
	var users []*User
	query := "username LIKE ? OR display_name LIKE ?"
	key := "%" + keyword + "%"
	err := tx.Where(query, key, key).Find(&users).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return users
}

func (*repositoryImpl) FindByEmail(tx *gorm.DB, email string) *User {
	var user *User
	err := tx.Where("email = ?", email).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

func (*repositoryImpl) FindByUsername(tx *gorm.DB, username string) *User {
	var user *User
	err := tx.Where("username = ?", username).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

func (*repositoryImpl) FindByEmailOrUsername(tx *gorm.DB, handler string) *User {
	var user *User
	err := tx.Where("email = ? OR username = ?", handler, handler).
		Limit(1).
		Find(&user).Error
	if err != nil {
		panic(exception.DatabaseError{Message: err.Error()})
	}
	return user
}

package app

import (
	"context"
	"go-api/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"time"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Init Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:root@tcp(localhost:3306)/go_api?parseTime=true&loc=Local",
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	helper.PanicIfError(err)

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	//db.LogMode(true)
	DB = db
	return DB
}

// TestDBInit This function will create a temporary database for running testing cases
func TestDBInit() *gorm.DB {
	testDB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:root@tcp(localhost:3306)/go_api_test?parseTime=true&loc=Local",
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	helper.PanicIfError(err)

	sqlDB, err := testDB.DB()
	sqlDB.SetMaxIdleConns(5)

	DB = testDB
	return DB
}

// GetDB Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}

func Tx(ctx context.Context, task func(tx *gorm.DB)) {
	tx := DB.WithContext(ctx)
	task(tx)
	defer callbacks.CommitOrRollbackTransaction(tx)
	//err := tx.Transaction(func(tx *gorm.DB) error {
	//	return task(tx)
	//})
	//helper.PanicIfError(err)
}

func InsertUsingTx(data ...interface{}) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		for _, v := range data {
			if err := tx.Create(v).Error; err != nil {
				return err
			}
		}
		return nil
	}
}

package helper

import (
	"gorm.io/gorm"
)

func TXCommitOrRollback(tx *gorm.DB) {
	err := recover()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}
}

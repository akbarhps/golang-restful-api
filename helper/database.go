package helper

import (
	"database/sql"
)

func TXCommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

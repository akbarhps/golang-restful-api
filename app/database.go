package app

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(env string) *sql.DB {
	source := "root:root@tcp(localhost:3306)/go_api_test?parseTime=true"
	if env == "prod" {
		source = "root:root@tcp(localhost:3306)/go_api?parseTime=true"
	}

	db, err := sql.Open("mysql", source)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

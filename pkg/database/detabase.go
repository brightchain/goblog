package database

import (
	"database/sql"
	"goblog/pkg/logger"
	"time"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Initialize() {
	var err error
	config := mysql.Config{
		User:   "root",
		Passwd: "ch123213",
		Addr:   "192.168.9.28:3306",
		Net:    "tcp",
		DBName: "goblog",
	}

	DB, err = sql.Open("mysql", config.FormatDSN())

	logger.LogError(err)

	DB.SetMaxOpenConns(100)
	DB.SetConnMaxIdleTime(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	err = DB.Ping()
	logger.LogError(err)
}

func CreateArticlesSQL() {
	sql := `
	     CREATE TABLE IF NOT EXISTS articles (
		 	id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
			title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
			body longtext COLLATE utf8mb4_unicode_ci NOT NULL
		 )
	`
	_, err := DB.Exec(sql)
	logger.LogError(err)
}

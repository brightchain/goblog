package database

import (
	"database/sql"
	"fmt"
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"time"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Initialize() {
	var err error
	config := mysql.Config{
		User:                 config.GetString("database.mysql.username"),
		Passwd:               config.GetString("database.mysql.password"),
		Addr:                 fmt.Sprintf("%v:%v", config.GetString("database.mysql.host"), config.GetString("database.mysql.port")),
		Net:                  "tcp",
		DBName:               config.GetString("database.mysql.username"),
		AllowNativePasswords: true,
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

package model

import (
	"goblog/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	config := mysql.New(mysql.Config{
		DSN: "root:ch123213@tcp(192.168.9.22:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	})
	DB, err = gorm.Open(config, &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Warn)})
	logger.LogError(err)

	return DB
}

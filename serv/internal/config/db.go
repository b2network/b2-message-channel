package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToDB(cfg AppConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Database.UserName, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)
	log.Info("connect to ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.Database.LogLevel)),
	})
	if err != nil {
		panic(err)
	}
	return db
}

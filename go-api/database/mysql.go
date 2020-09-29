package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/peashoot/peapi/config"
)

// Connect 连接
func Connect() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Config.DatabaseConfig.DBUsername, config.Config.DatabaseConfig.DBPassword, config.Config.DatabaseConfig.DBHost, config.Config.DatabaseConfig.DBName))
	if err != nil {
		log.Println("Failure to connect mysql databse, the detail of error is", err)
	}
	return db, err
}

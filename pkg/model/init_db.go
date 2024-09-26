package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := database.AutoMigrate(&Event{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}

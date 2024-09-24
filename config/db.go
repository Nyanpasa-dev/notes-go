package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"simple-api/models"
)

var db *gorm.DB

func RunDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		AppConfig.Database.Host,
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.DBName,
		AppConfig.Database.Port,
		AppConfig.Database.SSLMode,
		AppConfig.Database.TimeZone,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Note{}, &models.User{})
	if err != nil {
		panic("failed to migrate database")
	}
}
func GetDB() *gorm.DB {
	return db
}

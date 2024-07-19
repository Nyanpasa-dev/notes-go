package config

import (
	"simple-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	// Initialize the database
	dsn := "host=localhost user=postgres password=root dbname=notes-go port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Note{}, &models.User{})

}
func GetDB() *gorm.DB {
	return db
}

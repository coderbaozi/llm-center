package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/coderbaozi/llm-center/models"
)

var db *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Product{})
}

package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coderbaozi/llm-center/models"
	"gopkg.in/yaml.v3"
	"os"
)

var DB *gorm.DB

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

func InitDB() {
	config, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.Charset,
		config.Database.ParseTime,
		config.Database.Loc)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Product{}, &models.Agent{})
}

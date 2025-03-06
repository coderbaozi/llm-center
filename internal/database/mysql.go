package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()
	if dsn == "" {
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Driver)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	DB = db
	log.Println("数据库连接初始化成功")
	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}

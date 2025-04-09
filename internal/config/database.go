package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/llm-center/internal/model"
)

func (c *DatabaseConfig) MakeMysqlUrl() string {
	switch c.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
	default:
		return ""
	}
}

var DB *gorm.DB

func InitDB(cfg *DatabaseConfig) error {
	mysqlUrl := cfg.MakeMysqlUrl()
	if mysqlUrl == "" {
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Driver)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(mysqlUrl), gormConfig)
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

func GetDB() *gorm.DB {
	return DB
}

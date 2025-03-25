package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _DbMaster *gorm.DB
var _DbSlave *gorm.DB

// DB 返回只读实例
func DB() *gorm.DB {
	return _DbSlave
}

// DBMaster 返回主库实例
func DBMaster() *gorm.DB {
	return _DbMaster
}

func init() {
	_DbSlave, _ = NewDB(Cfg)
	_DbMaster, _ = NewDB(Cfg)
}

// NewDB 创建新的数据库连接
func NewDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(&Alarm{})
	if err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	return db, nil
}

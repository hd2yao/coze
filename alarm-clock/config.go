package main

// 加载全局配置

var Cfg = NewConfig()

// DBConfig 数据库配置
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Config 应用配置
type Config struct {
	DB DBConfig
}

// NewConfig 创建默认配置
func NewConfig() *Config {
	return &Config{
		DB: DBConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "root",
			DBName:   "coze",
		},
	}
}

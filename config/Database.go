package config

import (
	"fmt"
	"gorm.io/gorm"
)

var DB *gorm.DB

//DBConfig represents db configuration

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "tutorview_gouser",
		Password: "bQGgKwTdh9",
		DBName:   "tutorview_godb",
		//Host:     "localhost",
		//User:     "root",
		//Password: "",
		//DBName:   "tutorview_db",
	}
	return &dbConfig
}

func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}

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
	//port, err := strconv.Atoi(os.Getenv("APP_DB_PORT"))
	/*
		if err != nil {
			fmt.Println("Error getting port: ", err)
		}*/
	dbConfig := DBConfig{
		Port:     3306,
		Host:     "localhost",
		User:     "base",
		Password: "Serve@2022.Cloud",
		DBName:   "wholesome_ed_db",
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

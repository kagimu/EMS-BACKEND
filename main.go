//package tutorView
package main

import (
	"database/sql"
	"fmt"
	"github.com/hamsoftug/wholesomeEd/config"
	"github.com/hamsoftug/wholesomeEd/models"
	"github.com/hamsoftug/wholesomeEd/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

func main() {
	//config.DB, err = gorm.Open("mysql", config.DbURL(config.BuildDBConfig()))
	config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{})

	//config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	//All auto migrations
	//defer config.DB.Close()
	sqlDB, err := config.DB.DB()
	if err != nil {
		return
	}
	//All auto migrations
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {

		}
	}(sqlDB)

	errMigrate := config.DB.AutoMigrate(
		&models.User{},
		&models.Contact{},
		&models.PasswordReset{},
	)
	if errMigrate != nil {
		return
	}

	router := routes.SetupRoute()

	errPort := router.Run(":2053")
	if errPort != nil {
		fmt.Println("Error starting server:", errPort)
	}

}

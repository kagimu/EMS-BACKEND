//package tutorView
package main

import (
	"database/sql"
	"fmt"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/api/routes"
	models2 "github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"github.com/hamsoftug/wholesomeEd/domain/campus/models"
	models4 "github.com/hamsoftug/wholesomeEd/domain/class/models"
	models6 "github.com/hamsoftug/wholesomeEd/domain/learningarea/models"
	models3 "github.com/hamsoftug/wholesomeEd/domain/school/models"
	models5 "github.com/hamsoftug/wholesomeEd/domain/subject/models"
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
		&models2.User{},
		&models2.Contact{},
		&models2.PasswordReset{},
		&models2.UserPassword{},
		&models3.School{},
		&models.Campus{},
		&models.Level{},
		&models4.Class{},
		&models5.Subject{},
		&models6.LearningArea{},
		&models4.SubjectCLass{},
	)
	if errMigrate != nil {
		return
	}

	router := routes.SetupRoute()

	errPort := router.Run(":2054")
	if errPort != nil {
		fmt.Println("Error starting server:", errPort)
	}

}

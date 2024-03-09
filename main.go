package main

import (
	"database/sql"
	"fmt"

	"os"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/logger"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

func main() {
	//if os.Getenv("ENVIRONMENT") != "production" {
	//	err := godotenv.Load(".env")
	//	if err != nil {
	//		log.Fatal("Error loading .env file")
	//	}
	//}
	//
	//if os.Getenv("ENVIRONMENT") == "production" {
	//	gin.SetMode(gin.ReleaseMode)
	//}

	errZone := os.Setenv("TZ", "Africa/Nairobi")
	if errZone != nil {
		return
	}

	config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	sqlDB, err := config.DB.DB()
	if err != nil {
		return
	}

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			fmt.Println("Error closing database:", err)
		}
	}(sqlDB)

	//All auto migrations
	if err := migrations.Init(); err != nil {
		log.Println(err)
	}

	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	loger, err := logger.Init(rawJSON)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	loger.Info("Logger construction succeeded")

	router := routes.SetupRoutes()
	port := os.Getenv("APP_PORT")

	server := router.Run(":" + port)
	if server != nil {
		fmt.Println("Error starting server:", server)
	}
}

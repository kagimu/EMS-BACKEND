package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"log"
	"os"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	logs "github.com/Noregrets-Systems/wholesomeEdBackEnd/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	/*err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}*/

	gin.SetMode(gin.DebugMode)

	err := os.Setenv("TZ", "Africa/Nairobi")
	if err != nil {
		return
	}

	config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	sqlDB, err := config.DB.DB()
	if err != nil {
		fmt.Println("Error getting database connection:", err)
		return
	}

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			fmt.Println("Error closing database:", err)
			return
		}
	}(sqlDB)

	//All auto migrations
	if err := migrations.Init(); err != nil {
		log.Println(err)
		return
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

	logg, err := logs.Init(rawJSON)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	logg.Info("Logger construction succeeded")

	router := routes.SetupRoutes()
	port := os.Getenv("APP_PORT")

	server := router.Run(":" + port)
	if server != nil {
		fmt.Println("Error starting server:", server)
	}
}

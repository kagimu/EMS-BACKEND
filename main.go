package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	commercialServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize the timezone globally
func init() {
	loc, err := time.LoadLocation("Africa/Kampala")
	if err != nil {
		panic(fmt.Sprintf("Failed to load timezone: %v", err))
	}
	time.Local = loc // Set the default local timezone
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// EMS Release 1 (security fix B1): fail fast if the JWT signing secret is
	// absent. This replaces the removed "secret" fallback so that a predictable
	// signing key can never be used in any environment.
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set; refusing to start")
	}

	gin.SetMode(gin.DebugMode)

	err = os.Setenv("TZ", "Africa/Nairobi")
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

	// Development-only, idempotent illustrative package catalogue for the
	// disposable onboarding environment. Never enabled in production.
	if strings.EqualFold(strings.TrimSpace(os.Getenv("DEV_SEED_PACKAGES")), "true") {
		if n, err := commercialServices.SeedTestPackages(); err != nil {
			log.Println("test package seed failed:", err)
		} else if n > 0 {
			log.Printf("seeded %d illustrative installation package(s)", n)
		}
	}

	router := routes.SetupRoutes()
	port := os.Getenv("APP_PORT")

	server := router.Run(":" + port)
	if server != nil {
		fmt.Println("Error starting server:", server)
	}
}

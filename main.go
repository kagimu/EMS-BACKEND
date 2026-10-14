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
	// .env is optional: it is used for local development, but in the container
	// (e.g. Render) configuration is supplied through process environment
	// variables. A missing file must not prevent startup.
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found; using process environment variables")
	}

	// EMS Release 1 (security fix B1): fail fast if the JWT signing secret is
	// absent. This replaces the removed "secret" fallback so that a predictable
	// signing key can never be used in any environment.
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set; refusing to start")
	}

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

	// Prefer APP_PORT; fall back to the platform-provided PORT (Render sets
	// PORT), then a local default, so the container always binds a real port.
	port := os.Getenv("APP_PORT")
	if strings.TrimSpace(port) == "" {
		port = os.Getenv("PORT")
	}
	if strings.TrimSpace(port) == "" {
		port = "2053"
	}

	server := router.Run(":" + port)
	if server != nil {
		fmt.Println("Error starting server:", server)
	}
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	commercialServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/services"
	paymentServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ..
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

	// Apply the schema migration WITHOUT blocking the listening port.
	//
	// Render (and similar platforms) fail the deploy if the container does not
	// accept connections on PORT promptly, and GORM AutoMigrate against a remote
	// database can take minutes (it issues several introspection queries per
	// table at ~150 ms each over the network). The HTTP server therefore opens
	// the port immediately, and migrations (plus the optional dev seed) run in
	// the background. On a first deploy, requests that touch a not-yet-created
	// table may fail until the migration finishes; subsequent deploys find the
	// schema already present.
	workerCtx, stopWorker := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopWorker()
	go func() {
		if err := migrations.Init(); err != nil {
			log.Println("schema migration error:", err)
			return
		}
		// Development-only, idempotent illustrative package catalogue for the
		// disposable onboarding environment. Never enabled in production.
		if commercialServices.SeedingEnabled() {
			if n, err := commercialServices.SeedTestPackages(); err != nil {
				log.Println("test package seed failed:", err)
			} else if n > 0 {
				log.Printf("seeded %d illustrative installation package(s)", n)
			}
		}
		paymentServices.StartReconciliationWorker(workerCtx)
	}()

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

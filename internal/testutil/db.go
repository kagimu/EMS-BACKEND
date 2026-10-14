// Package testutil provides shared helpers for the EMS Release 1 regression
// foundation.
//
// It contains no production code and is imported only by tests. It exists so
// that integration tests can connect to a real MySQL/MariaDB test database
// without duplicating connection logic, while remaining skippable in
// environments that do not have a database configured.
package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// RequireTestDB connects to the integration test database described by the
// TEST_APP_DB_* environment variables and assigns it to config.DB so that the
// repositories and handlers under test use it.
//
// When TEST_APP_DB_HOST is unset the calling test is skipped, which keeps
// `go test ./...` runnable in environments without a database.
func RequireTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	host := os.Getenv("TEST_APP_DB_HOST")
	if host == "" {
		t.Skip("TEST_APP_DB_HOST not set; skipping integration test")
	}

	user := os.Getenv("TEST_APP_DB_USERNAME")
	pass := os.Getenv("TEST_APP_DB_PASSWORD")
	port := os.Getenv("TEST_APP_DB_PORT")
	name := os.Getenv("TEST_APP_DB_NAME")
	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("connecting to test database: %v", err)
	}

	config.DB = db

	return db
}

// LoginCredentials returns the email and password of the seeded integration
// test user, or skips the test when they are not configured.
//
// Release 1 deliberately does not invent a fixture schema. The test database
// is expected to be seeded by the team with a known administrator, and the
// credentials are supplied through the environment.
func LoginCredentials(t *testing.T) (email string, password string) {
	t.Helper()

	email = os.Getenv("TEST_LOGIN_EMAIL")
	password = os.Getenv("TEST_LOGIN_PASSWORD")
	if email == "" || password == "" {
		t.Skip("TEST_LOGIN_EMAIL/TEST_LOGIN_PASSWORD not set; skipping authenticated flow")
	}

	return email, password
}

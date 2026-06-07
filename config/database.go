package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection instance.
var DB *gorm.DB

// ConnectDB establishes a connection to the PostgreSQL database using
// environment variables and assigns it to the global DB variable.
// It also runs AutoMigrate for all registered models.
func ConnectDB(models ...interface{}) {
	dsn := buildDSN()

	gormConfig := &gorm.Config{
		Logger: buildLogger(),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("[DB] Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[DB] Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("[DB] Database connection established successfully.")

	// Run AutoMigrate for all provided model structs
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			log.Fatalf("[DB] AutoMigrate failed: %v", err)
		}
		log.Println("[DB] AutoMigrate completed successfully.")
	}

	DB = db
}

// buildDSN constructs the PostgreSQL DSN string from environment variables.
func buildDSN() string {
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "glk_user")
	password := GetEnv("DB_PASSWORD", "glk_secret")
	dbname := GetEnv("DB_NAME", "glk_db")
	sslmode := GetEnv("DB_SSLMODE", "disable")
	timezone := GetEnv("DB_TIMEZONE", "Asia/Jakarta")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbname, sslmode, timezone,
	)
}

// buildLogger returns a GORM logger configured based on the APP_ENV variable.
// In production it only logs errors; in development it logs all queries.
func buildLogger() logger.Interface {
	appEnv := GetEnv("APP_ENV", "development")

	logLevel := logger.Info
	if appEnv == "production" {
		logLevel = logger.Error
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  appEnv != "production",
		},
	)
}

// GetEnv retrieves an environment variable value or returns a fallback default.
// Exported so other packages (e.g. main) can use it.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

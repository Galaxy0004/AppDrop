// Package database manages the lifecycle of the application's connection to the persistent data store.
// It handles configuration parsing, connection establishment, and pooling parameters.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DB represents the global database connection pool instance.
var DB *sql.DB

// Config encapsulates the necessary parameters for establishing a connection to a PostgreSQL database.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfigFromEnv initializes a Config instance by reading parameters from system environment variables.
// It applies sensible defaults if specific variables are not defined.
func NewConfigFromEnv() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "appdrop"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv retrieves the value of an environment variable or returns a specified default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Connect establishes a persistent connection to the PostgreSQL database using the provided configuration.
// It performs a ping to verify connectivity and configures connection pool parameters.
func Connect(config Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("Successfully connected to database")
	return nil
}

// Close terminates the active database connection pool if it exists.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB retrieves the active global database connection pool.
func GetDB() *sql.DB {
	return DB
}

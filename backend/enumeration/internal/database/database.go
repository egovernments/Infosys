// Package database handles database connection and access
package database

import (
	"enumeration/internal/config"
	"enumeration/pkg/logger"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// dbInstance holds the singleton GORM DB instance for the application.
// dbOnce ensures the database is initialized only once (thread-safe singleton pattern).
var (
	dbInstance *gorm.DB  // Singleton DB instance
	dbOnce     sync.Once // Ensures DB is initialized once
)

// Connect creates a new database connection using the provided config.
// It sets up the GORM ORM with a custom schema naming strategy and returns the DB instance.
func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable client_encoding=UTF8 search_path=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSchema)

	// Initialize GORM with schema naming strategy
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "DIGIT3.", // Prefix for all table names
		},
	})
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	logger.Info("Database connected and migrated successfully")
	return db
}

// GetDB returns the singleton database instance.
// Initializes the DB connection on first call using sync.Once for thread safety.
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		dbInstance = Connect(config.GetConfig())
	})
	return dbInstance
}

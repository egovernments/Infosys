// Package config handles application configuration loading
package config

import (
	"os"
	"sync"
)

// Config holds all environment-based configuration values for the application.
type Config struct {
	DBHost      string // Database host
	DBPort      string // Database port
	DBUser      string // Database username
	DBPassword  string // Database password
	DBName      string // Database name
	DBSchema    string // Database schema
	Port        string // Service port for the application
	WorkflowURL string // Workflow service URL
	AuditURL    string // Audit service URL
	IDGenURL    string // ID generation service URL
	MDMSURL     string // MDMS service URL
}

// configInstance holds the singleton Config instance.
// configOnce ensures the config is loaded only once (thread-safe singleton pattern).
var (
	configInstance *Config   // Singleton instance
	configOnce     sync.Once // Ensures config is loaded once
)

// Load reads environment variables and returns a Config struct populated with values.
// If an environment variable is not set, a default value is used.
func Load() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "property"),
		DBSchema:    getEnv("DB_SCHEMA", "DIGIT3"),
		Port:        getEnv("PORT", "8080"),
		WorkflowURL: getEnv("WORKFLOW_URL", "http://localhost:8081"),
		AuditURL:    getEnv("AUDIT_URL", "http://localhost:8082/audit/logs"),
		IDGenURL:    getEnv("IDGEN_URL", "http://localhost:8083/idgen/v1/generate"),
		MDMSURL:     getEnv("MDMS_URL", "http://localhost:8084"),
	}
}

// GetConfig returns the singleton Config instance.
// Initializes the config on first call using sync.Once for thread safety.
func GetConfig() *Config {
	configOnce.Do(func() {
		configInstance = Load()
	})
	return configInstance
}

// getEnv returns the value of an environment variable or a default if not set.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// RoleConfig holds role-based permissions loaded from a YAML file.
type RoleConfig struct {
	RolePermissions map[string][]string `yaml:"rolePermissions"` // Map of role to permissions
}

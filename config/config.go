package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// AuthConfig contains authentication-specific configuration
type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret"`
}

// DatabaseConfig contains database-specific configuration
type DatabaseConfig struct {
	Type     string `yaml:"type"` // sqlite, postgres, mysql
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"` // For PostgreSQL: disable, require, verify-ca, verify-full
	Path     string `yaml:"path"`     // For SQLite: file path
	
	// Connection pool settings
	MaxOpenConns int `yaml:"max_open_conns"`
	MaxIdleConns int `yaml:"max_idle_conns"`
}

// LoadConfig loads configuration from a YAML file with environment variable overrides
func LoadConfig(configPath string) (*Config, error) {
	// Set default configuration
	config := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Type:         "sqlite",
			Path:         "./playground.db",
			MaxOpenConns: 100,
			MaxIdleConns: 10,
			SSLMode:      "disable",
		},
		Auth: AuthConfig{
			JWTSecret: "default-secret-change-in-production",
		},
	}

	// If config file exists, load it
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
			// File doesn't exist, use defaults
		} else {
			if err := yaml.Unmarshal(data, config); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		}
	}

	// Override with environment variables
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}
	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		config.Database.Type = dbType
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			config.Database.Port = p
		}
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.Name = dbName
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		config.Database.Path = dbPath
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.Auth.JWTSecret = jwtSecret
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate database type
	validTypes := map[string]bool{
		"sqlite":   true,
		"postgres": true,
		"mysql":    true,
	}
	if !validTypes[c.Database.Type] {
		return fmt.Errorf("invalid database type: %s (must be sqlite, postgres, or mysql)", c.Database.Type)
	}

	// Validate database-specific fields
	switch c.Database.Type {
	case "sqlite":
		if c.Database.Path == "" {
			return fmt.Errorf("database path is required for SQLite")
		}
	case "postgres", "mysql":
		if c.Database.Host == "" {
			return fmt.Errorf("database host is required for %s", c.Database.Type)
		}
		if c.Database.Port == 0 {
			// Set default ports
			if c.Database.Type == "postgres" {
				c.Database.Port = 5432
			} else if c.Database.Type == "mysql" {
				c.Database.Port = 3306
			}
		}
		if c.Database.Name == "" {
			return fmt.Errorf("database name is required for %s", c.Database.Type)
		}
		if c.Database.User == "" {
			return fmt.Errorf("database user is required for %s", c.Database.Type)
		}
	}

	return nil
}

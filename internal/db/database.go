package db

import (
	"fmt"
	"log"

	"github.com/Creeper19472/playground/config"
	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection using GORM with the specified configuration
func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	var dialector gorm.Dialector
	
	// Select the appropriate database driver based on configuration
	switch cfg.Type {
	case "sqlite":
		dialector = sqlite.Open(cfg.Path)
		log.Printf("Connecting to SQLite database: %s", cfg.Path)
		
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
		dialector = postgres.Open(dsn)
		log.Printf("Connecting to PostgreSQL database: %s@%s:%d/%s", cfg.User, cfg.Host, cfg.Port, cfg.Name)
		
	case "mysql":
		// DSN format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		dialector = mysql.Open(dsn)
		log.Printf("Connecting to MySQL database: %s@%s:%d/%s", cfg.User, cfg.Host, cfg.Port, cfg.Name)
		
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// PrepareStmt improves performance
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Get generic database object to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings for high concurrency
	maxOpenConns := cfg.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 100
	}
	maxIdleConns := cfg.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 10
	}
	
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return &Database{DB: db}, nil
}

// NewDatabaseLegacy creates a new database connection using the legacy method (for backward compatibility)
func NewDatabaseLegacy(dbPath string) (*Database, error) {
	cfg := &config.DatabaseConfig{
		Type:         "sqlite",
		Path:         dbPath,
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}
	return NewDatabase(cfg)
}

// AutoMigrate runs automatic migration for all models
func (d *Database) AutoMigrate() error {
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.Issue{},
		&models.Vote{},
		&models.Message{},
		&models.Reference{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database schema migrated successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

const (
	// CredentialsTimestampFormat is the format used for timestamps in the credentials file
	CredentialsTimestampFormat = "2006-01-02 15:04:05 MST"
)

// AdminInitializer handles the creation of the initial admin account
type AdminInitializer struct {
	db *gorm.DB
}

// NewAdminInitializer creates a new admin initializer
func NewAdminInitializer(db *gorm.DB) *AdminInitializer {
	return &AdminInitializer{db: db}
}

// GenerateSecurePassword generates a cryptographically secure random password
func GenerateSecurePassword(length int) (string, error) {
	// Generate random bytes
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	// Encode to base64 for readable ASCII characters
	// Use URL encoding to avoid special characters that might cause issues
	password := base64.URLEncoding.EncodeToString(bytes)
	
	// Trim to desired length
	if len(password) > length {
		password = password[:length]
	}
	
	return password, nil
}

// InitializeAdminAccount checks if an admin account exists and creates one if not
// Returns (created bool, username string, password string, error)
func (ai *AdminInitializer) InitializeAdminAccount() (bool, string, string, error) {
	// Check if any admin user exists
	var count int64
	if err := ai.db.Model(&models.User{}).Where("is_admin = ?", true).Count(&count).Error; err != nil {
		return false, "", "", fmt.Errorf("failed to check for existing admin: %w", err)
	}
	
	// If admin already exists, do nothing
	if count > 0 {
		return false, "", "", nil
	}
	
	// Generate secure credentials
	username := "admin"
	password, err := GenerateSecurePassword(24)
	if err != nil {
		return false, "", "", fmt.Errorf("failed to generate secure password: %w", err)
	}
	
	// Create admin user
	admin := &models.User{
		Username: username,
		Email:    "admin@localhost",
		IsAdmin:  true,
		IsActive: true,
	}
	
	if err := admin.SetPassword(password); err != nil {
		return false, "", "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	if err := ai.db.Create(admin).Error; err != nil {
		return false, "", "", fmt.Errorf("failed to create admin user: %w", err)
	}
	
	return true, username, password, nil
}

// WriteCredentialsToFile securely writes admin credentials to a file
func WriteCredentialsToFile(username, password, filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Prepare credentials content
	content := fmt.Sprintf(`ADMINISTRATOR CREDENTIALS
========================
Generated: %s

Username: %s
Password: %s

IMPORTANT SECURITY NOTES:
-------------------------
1. This file contains sensitive credentials. Keep it secure!
2. Delete this file after you've saved the credentials elsewhere.
3. Change the admin password immediately after first login.
4. Do not commit this file to version control.
5. Restrict file permissions (0600) if not already set.

To change the password after login:
POST /api/v1/auth/change-password
{
  "current_password": "<current>",
  "new_password": "<new>",
  "confirm_password": "<new>"
}

`, time.Now().Format(CredentialsTimestampFormat), username, password)
	
	// Write file with restrictive permissions (owner read/write only)
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}
	
	return nil
}

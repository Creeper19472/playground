package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Creeper19472/playground/internal/models"
)

func TestGenerateSecurePassword(t *testing.T) {
	password, err := GenerateSecurePassword(24)
	if err != nil {
		t.Fatalf("Failed to generate password: %v", err)
	}

	if len(password) != 24 {
		t.Errorf("Expected password length 24, got %d", len(password))
	}

	// Test that multiple calls generate different passwords
	password2, err := GenerateSecurePassword(24)
	if err != nil {
		t.Fatalf("Failed to generate second password: %v", err)
	}

	if password == password2 {
		t.Error("Expected different passwords, got same")
	}
}

func TestInitializeAdminAccount_FirstRun(t *testing.T) {
	db := setupTestDB(t)
	adminInit := NewAdminInitializer(db)

	created, username, password, err := adminInit.InitializeAdminAccount()
	if err != nil {
		t.Fatalf("Failed to initialize admin account: %v", err)
	}

	if !created {
		t.Error("Expected admin account to be created")
	}

	if username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", username)
	}

	if len(password) != 24 {
		t.Errorf("Expected password length 24, got %d", len(password))
	}

	// Verify admin was created in database
	var user models.User
	if err := db.Where("username = ?", "admin").First(&user).Error; err != nil {
		t.Fatalf("Failed to find created admin: %v", err)
	}

	if !user.IsAdmin {
		t.Error("Expected user to be admin")
	}

	if !user.IsActive {
		t.Error("Expected user to be active")
	}

	if user.Email != "admin@localhost" {
		t.Errorf("Expected email 'admin@localhost', got '%s'", user.Email)
	}

	// Verify password works
	if !user.CheckPassword(password) {
		t.Error("Password verification failed")
	}
}

func TestInitializeAdminAccount_AlreadyExists(t *testing.T) {
	db := setupTestDB(t)
	adminInit := NewAdminInitializer(db)

	// Create initial admin
	created, _, _, err := adminInit.InitializeAdminAccount()
	if err != nil {
		t.Fatalf("Failed to initialize admin account: %v", err)
	}

	if !created {
		t.Error("Expected first admin to be created")
	}

	// Try to create again
	created2, username2, password2, err := adminInit.InitializeAdminAccount()
	if err != nil {
		t.Fatalf("Failed on second initialization: %v", err)
	}

	if created2 {
		t.Error("Expected admin not to be created on second run")
	}

	if username2 != "" {
		t.Error("Expected empty username when admin already exists")
	}

	if password2 != "" {
		t.Error("Expected empty password when admin already exists")
	}
}

func TestWriteCredentialsToFile(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_credentials.txt")

	username := "testadmin"
	password := "testpassword123"

	err := WriteCredentialsToFile(username, password, filePath)
	if err != nil {
		t.Fatalf("Failed to write credentials file: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("Credentials file was not created")
	}

	// Check file permissions (should be 0600)
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat credentials file: %v", err)
	}

	mode := info.Mode()
	if mode.Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", mode.Perm())
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read credentials file: %v", err)
	}

	contentStr := string(content)
	if len(contentStr) == 0 {
		t.Error("Credentials file is empty")
	}

	// Check that username and password are in the file
	if !contains(contentStr, username) {
		t.Error("Username not found in credentials file")
	}

	if !contains(contentStr, password) {
		t.Error("Password not found in credentials file")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

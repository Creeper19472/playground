package services

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create schema
	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

func TestUserService_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	user, err := service.CreateUser("testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	// Create a user
	created, err := service.CreateUser("testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user
	user, err := service.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.ID != created.ID {
		t.Errorf("Expected user ID %d, got %d", created.ID, user.ID)
	}

	if user.Username != created.Username {
		t.Errorf("Expected username '%s', got '%s'", created.Username, user.Username)
	}
}

func TestUserService_ListUsers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	// Create multiple users
	_, err := service.CreateUser("user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	_, err = service.CreateUser("user2", "user2@example.com")
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// List users
	users, err := service.ListUsers()
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}

package services

import (
	"testing"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIssueTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.Issue{}, &models.Vote{}); err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Insert test users
	db.Create(&models.User{Username: "testuser", Email: "test@example.com"})
	db.Create(&models.User{Username: "voter", Email: "voter@example.com"})

	return db
}

func TestIssueService_CreateIssue(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	issue, err := service.CreateIssue(1, "issue", "Test Issue", "This is a test")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	if issue.ID == 0 {
		t.Error("Expected issue ID to be set")
	}

	if issue.Summary != "Test Issue" {
		t.Errorf("Expected summary 'Test Issue', got '%s'", issue.Summary)
	}

	if issue.Type != "issue" {
		t.Errorf("Expected type 'issue', got '%s'", issue.Type)
	}

	if issue.VoteCount != 0 {
		t.Errorf("Expected vote count 0, got %d", issue.VoteCount)
	}
}

func TestIssueService_VoteIssue(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	// Create an issue
	issue, err := service.CreateIssue(1, "issue", "Test Issue", "This is a test")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	// Vote on the issue
	err = service.VoteIssue(2, issue.ID)
	if err != nil {
		t.Fatalf("Failed to vote on issue: %v", err)
	}

	// Check vote count
	updatedIssue, err := service.GetIssueByID(issue.ID)
	if err != nil {
		t.Fatalf("Failed to get issue: %v", err)
	}

	if updatedIssue.VoteCount != 1 {
		t.Errorf("Expected vote count 1, got %d", updatedIssue.VoteCount)
	}

	// Try to vote again (should fail)
	err = service.VoteIssue(2, issue.ID)
	if err == nil {
		t.Error("Expected error when voting twice, got nil")
	}
}

func TestIssueService_UnvoteIssue(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	// Create an issue
	issue, err := service.CreateIssue(1, "issue", "Test Issue", "This is a test")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	// Vote on the issue
	err = service.VoteIssue(2, issue.ID)
	if err != nil {
		t.Fatalf("Failed to vote on issue: %v", err)
	}

	// Unvote
	err = service.UnvoteIssue(2, issue.ID)
	if err != nil {
		t.Fatalf("Failed to unvote issue: %v", err)
	}

	// Check vote count
	updatedIssue, err := service.GetIssueByID(issue.ID)
	if err != nil {
		t.Fatalf("Failed to get issue: %v", err)
	}

	if updatedIssue.VoteCount != 0 {
		t.Errorf("Expected vote count 0, got %d", updatedIssue.VoteCount)
	}
}

func TestIssueService_ListIssues(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	// Create multiple issues
	issue1, _ := service.CreateIssue(1, "issue", "Issue 1", "First issue")
	issue2, _ := service.CreateIssue(1, "issue", "Issue 2", "Second issue")

	// Vote on issue 2 to make it rank higher
	service.VoteIssue(2, issue2.ID)

	// List issues
	issues, err := service.ListIssues()
	if err != nil {
		t.Fatalf("Failed to list issues: %v", err)
	}

	if len(issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(issues))
	}

	// Check that issues are sorted by vote count (issue2 should be first)
	if issues[0].ID != issue2.ID {
		t.Errorf("Expected first issue to be issue2 (ID %d), got ID %d", issue2.ID, issues[0].ID)
	}

	if issues[1].ID != issue1.ID {
		t.Errorf("Expected second issue to be issue1 (ID %d), got ID %d", issue1.ID, issues[1].ID)
	}
}

func TestIssueService_CreateProposition(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	issue, err := service.CreateIssue(1, "proposition", "The earth is round", "Scientific fact")
	if err != nil {
		t.Fatalf("Failed to create proposition: %v", err)
	}

	if issue.Type != "proposition" {
		t.Errorf("Expected type 'proposition', got '%s'", issue.Type)
	}
}

func TestIssueService_InvalidType(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	_, err := service.CreateIssue(1, "invalid", "Bad type", "Should fail")
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

func TestIssueService_UpdateWithConclusion(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	issue, _ := service.CreateIssue(1, "issue", "Test Issue", "Description")

	updated, err := service.UpdateIssue(issue.ID, "", "", "The conclusion is X", nil)
	if err != nil {
		t.Fatalf("Failed to update issue: %v", err)
	}

	if updated.Conclusion != "The conclusion is X" {
		t.Errorf("Expected conclusion 'The conclusion is X', got '%s'", updated.Conclusion)
	}
}

func TestIssueService_UpdateWithTruthValue(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	issue, _ := service.CreateIssue(1, "proposition", "The earth is round", "Scientific fact")

	truthValue := true
	updated, err := service.UpdateIssue(issue.ID, "", "", "", &truthValue)
	if err != nil {
		t.Fatalf("Failed to update issue: %v", err)
	}

	if updated.TruthValue == nil || *updated.TruthValue != true {
		t.Error("Expected truth value to be true")
	}
}

func TestIssueService_DefaultType(t *testing.T) {
	db := setupIssueTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewIssueService(db)

	// Empty type should default to "issue"
	issue, err := service.CreateIssue(1, "", "Default type", "Should be issue")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	if issue.Type != "issue" {
		t.Errorf("Expected default type 'issue', got '%s'", issue.Type)
	}
}

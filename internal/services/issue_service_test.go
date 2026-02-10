package services

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupIssueTestDB(t *testing.T) *sql.DB {
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

	CREATE TABLE issues (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		summary TEXT NOT NULL,
		description TEXT,
		vote_count INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		issue_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, issue_id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (issue_id) REFERENCES issues(id)
	);

	INSERT INTO users (username, email) VALUES ('testuser', 'test@example.com');
	INSERT INTO users (username, email) VALUES ('voter', 'voter@example.com');
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

func TestIssueService_CreateIssue(t *testing.T) {
	db := setupIssueTestDB(t)
	defer db.Close()

	service := NewIssueService(db)

	issue, err := service.CreateIssue(1, "Test Issue", "This is a test")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	if issue.ID == 0 {
		t.Error("Expected issue ID to be set")
	}

	if issue.Summary != "Test Issue" {
		t.Errorf("Expected summary 'Test Issue', got '%s'", issue.Summary)
	}

	if issue.VoteCount != 0 {
		t.Errorf("Expected vote count 0, got %d", issue.VoteCount)
	}
}

func TestIssueService_VoteIssue(t *testing.T) {
	db := setupIssueTestDB(t)
	defer db.Close()

	service := NewIssueService(db)

	// Create an issue
	issue, err := service.CreateIssue(1, "Test Issue", "This is a test")
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
	defer db.Close()

	service := NewIssueService(db)

	// Create an issue
	issue, err := service.CreateIssue(1, "Test Issue", "This is a test")
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
	defer db.Close()

	service := NewIssueService(db)

	// Create multiple issues
	issue1, _ := service.CreateIssue(1, "Issue 1", "First issue")
	issue2, _ := service.CreateIssue(1, "Issue 2", "Second issue")

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

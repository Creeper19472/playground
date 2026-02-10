package services

import (
	"database/sql"
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
)

type IssueService struct {
	db *sql.DB
}

func NewIssueService(db *sql.DB) *IssueService {
	return &IssueService{db: db}
}

// CreateIssue creates a new issue
func (s *IssueService) CreateIssue(userID int64, summary, description string) (*models.Issue, error) {
	result, err := s.db.Exec(
		"INSERT INTO issues (user_id, summary, description) VALUES (?, ?, ?)",
		userID, summary, description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get issue ID: %w", err)
	}

	return s.GetIssueByID(id)
}

// GetIssueByID retrieves an issue by ID
func (s *IssueService) GetIssueByID(id int64) (*models.Issue, error) {
	issue := &models.Issue{}
	err := s.db.QueryRow(`
		SELECT i.id, i.user_id, u.username, i.summary, i.description, i.vote_count, i.created_at, i.updated_at
		FROM issues i
		JOIN users u ON i.user_id = u.id
		WHERE i.id = ?`,
		id,
	).Scan(&issue.ID, &issue.UserID, &issue.Username, &issue.Summary, &issue.Description, &issue.VoteCount, &issue.CreatedAt, &issue.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	return issue, nil
}

// ListIssues retrieves all issues, ordered by vote count
func (s *IssueService) ListIssues() ([]*models.Issue, error) {
	rows, err := s.db.Query(`
		SELECT i.id, i.user_id, u.username, i.summary, i.description, i.vote_count, i.created_at, i.updated_at
		FROM issues i
		JOIN users u ON i.user_id = u.id
		ORDER BY i.vote_count DESC, i.created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}
	defer rows.Close()

	var issues []*models.Issue
	for rows.Next() {
		issue := &models.Issue{}
		if err := rows.Scan(&issue.ID, &issue.UserID, &issue.Username, &issue.Summary, &issue.Description, &issue.VoteCount, &issue.CreatedAt, &issue.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan issue: %w", err)
		}
		issues = append(issues, issue)
	}

	return issues, nil
}

// VoteIssue allows a user to vote on an issue
func (s *IssueService) VoteIssue(userID, issueID int64) error {
	// Check if user already voted
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM votes WHERE user_id = ? AND issue_id = ?", userID, issueID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing vote: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("user already voted on this issue")
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert vote
	_, err = tx.Exec("INSERT INTO votes (user_id, issue_id) VALUES (?, ?)", userID, issueID)
	if err != nil {
		return fmt.Errorf("failed to insert vote: %w", err)
	}

	// Update vote count
	_, err = tx.Exec("UPDATE issues SET vote_count = vote_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", issueID)
	if err != nil {
		return fmt.Errorf("failed to update vote count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UnvoteIssue allows a user to remove their vote from an issue
func (s *IssueService) UnvoteIssue(userID, issueID int64) error {
	// Check if user has voted
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM votes WHERE user_id = ? AND issue_id = ?", userID, issueID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing vote: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("user has not voted on this issue")
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete vote
	_, err = tx.Exec("DELETE FROM votes WHERE user_id = ? AND issue_id = ?", userID, issueID)
	if err != nil {
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	// Update vote count
	_, err = tx.Exec("UPDATE issues SET vote_count = vote_count - 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", issueID)
	if err != nil {
		return fmt.Errorf("failed to update vote count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUserVote checks if a user has voted on an issue
func (s *IssueService) GetUserVote(userID, issueID int64) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM votes WHERE user_id = ? AND issue_id = ?", userID, issueID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user vote: %w", err)
	}
	return count > 0, nil
}

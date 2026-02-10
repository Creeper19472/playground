package services

import (
	"database/sql"
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
)

type ReferenceService struct {
	db *sql.DB
}

func NewReferenceService(db *sql.DB) *ReferenceService {
	return &ReferenceService{db: db}
}

// CreateReference creates a new reference
func (s *ReferenceService) CreateReference(userID, issueID int64, url, title, description string) (*models.Reference, error) {
	result, err := s.db.Exec(
		"INSERT INTO references (user_id, issue_id, url, title, description) VALUES (?, ?, ?, ?, ?)",
		userID, issueID, url, title, description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create reference: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get reference ID: %w", err)
	}

	return s.GetReferenceByID(id)
}

// GetReferenceByID retrieves a reference by ID
func (s *ReferenceService) GetReferenceByID(id int64) (*models.Reference, error) {
	reference := &models.Reference{}
	err := s.db.QueryRow(
		"SELECT id, user_id, issue_id, url, title, description, created_at FROM references WHERE id = ?",
		id,
	).Scan(&reference.ID, &reference.UserID, &reference.IssueID, &reference.URL, &reference.Title, &reference.Description, &reference.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get reference: %w", err)
	}
	return reference, nil
}

// ListReferencesByIssue retrieves all references for an issue
func (s *ReferenceService) ListReferencesByIssue(issueID int64) ([]*models.Reference, error) {
	rows, err := s.db.Query(
		"SELECT id, user_id, issue_id, url, title, description, created_at FROM references WHERE issue_id = ? ORDER BY created_at DESC",
		issueID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list references: %w", err)
	}
	defer rows.Close()

	var references []*models.Reference
	for rows.Next() {
		reference := &models.Reference{}
		if err := rows.Scan(&reference.ID, &reference.UserID, &reference.IssueID, &reference.URL, &reference.Title, &reference.Description, &reference.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan reference: %w", err)
		}
		references = append(references, reference)
	}

	return references, nil
}

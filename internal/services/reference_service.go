package services

import (
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

type ReferenceService struct {
	db *gorm.DB
}

func NewReferenceService(db *gorm.DB) *ReferenceService {
	return &ReferenceService{db: db}
}

// CreateReference creates a new reference
func (s *ReferenceService) CreateReference(userID, issueID uint, url, title, description string) (*models.Reference, error) {
	reference := &models.Reference{
		UserID:      userID,
		IssueID:     issueID,
		URL:         url,
		Title:       title,
		Description: description,
	}
	
	if err := s.db.Create(reference).Error; err != nil {
		return nil, fmt.Errorf("failed to create reference: %w", err)
	}

	return reference, nil
}

// GetReferenceByID retrieves a reference by ID
func (s *ReferenceService) GetReferenceByID(id uint) (*models.Reference, error) {
	var reference models.Reference
	if err := s.db.First(&reference, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get reference: %w", err)
	}
	return &reference, nil
}

// ListReferencesByIssue retrieves all references for an issue
func (s *ReferenceService) ListReferencesByIssue(issueID uint) ([]*models.Reference, error) {
	var references []*models.Reference
	if err := s.db.Where("issue_id = ?", issueID).
		Order("created_at DESC").
		Find(&references).Error; err != nil {
		return nil, fmt.Errorf("failed to list references: %w", err)
	}
	return references, nil
}

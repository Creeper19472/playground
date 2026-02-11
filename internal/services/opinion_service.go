package services

import (
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

type OpinionService struct {
	db *gorm.DB
}

func NewOpinionService(db *gorm.DB) *OpinionService {
	return &OpinionService{db: db}
}

// CreateOpinion creates a new opinion under an issue
func (s *OpinionService) CreateOpinion(userID, issueID uint, content string) (*models.Opinion, error) {
	opinion := &models.Opinion{
		UserID:  userID,
		IssueID: issueID,
		Content: content,
	}

	if err := s.db.Create(opinion).Error; err != nil {
		return nil, fmt.Errorf("failed to create opinion: %w", err)
	}

	return s.GetOpinionByID(opinion.ID)
}

// GetOpinionByID retrieves an opinion by ID with user information
func (s *OpinionService) GetOpinionByID(id uint) (*models.Opinion, error) {
	var opinion models.Opinion
	if err := s.db.Preload("User").First(&opinion, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get opinion: %w", err)
	}

	opinion.Username = opinion.User.Username
	return &opinion, nil
}

// ListOpinionsByIssue retrieves all opinions for an issue
func (s *OpinionService) ListOpinionsByIssue(issueID uint) ([]*models.Opinion, error) {
	var opinions []*models.Opinion
	if err := s.db.Preload("User").
		Where("issue_id = ?", issueID).
		Order("created_at DESC").
		Find(&opinions).Error; err != nil {
		return nil, fmt.Errorf("failed to list opinions: %w", err)
	}

	for _, opinion := range opinions {
		opinion.Username = opinion.User.Username
	}

	return opinions, nil
}

// UpdateOpinion updates an opinion's content
func (s *OpinionService) UpdateOpinion(id uint, content string) (*models.Opinion, error) {
	var opinion models.Opinion
	if err := s.db.First(&opinion, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find opinion: %w", err)
	}

	if content != "" {
		opinion.Content = content
	}

	if err := s.db.Save(&opinion).Error; err != nil {
		return nil, fmt.Errorf("failed to update opinion: %w", err)
	}

	return s.GetOpinionByID(id)
}

// DeleteOpinion soft deletes an opinion
func (s *OpinionService) DeleteOpinion(id uint) error {
	if err := s.db.Delete(&models.Opinion{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete opinion: %w", err)
	}
	return nil
}

// AddSupport creates a support relationship: supporterID supports supportedID
func (s *OpinionService) AddSupport(supporterID, supportedID uint) error {
	var supporter, supported models.Opinion
	if err := s.db.First(&supporter, supporterID).Error; err != nil {
		return fmt.Errorf("supporter opinion not found: %w", err)
	}
	if err := s.db.First(&supported, supportedID).Error; err != nil {
		return fmt.Errorf("supported opinion not found: %w", err)
	}

	if err := s.db.Model(&supporter).Association("Supports").Append(&supported); err != nil {
		return fmt.Errorf("failed to add support relationship: %w", err)
	}

	return nil
}

// RemoveSupport removes a support relationship
func (s *OpinionService) RemoveSupport(supporterID, supportedID uint) error {
	var supporter, supported models.Opinion
	if err := s.db.First(&supporter, supporterID).Error; err != nil {
		return fmt.Errorf("supporter opinion not found: %w", err)
	}
	if err := s.db.First(&supported, supportedID).Error; err != nil {
		return fmt.Errorf("supported opinion not found: %w", err)
	}

	if err := s.db.Model(&supporter).Association("Supports").Delete(&supported); err != nil {
		return fmt.Errorf("failed to remove support relationship: %w", err)
	}

	return nil
}

// GetSupports retrieves all opinions that a given opinion supports
func (s *OpinionService) GetSupports(opinionID uint) ([]*models.Opinion, error) {
	var opinion models.Opinion
	if err := s.db.First(&opinion, opinionID).Error; err != nil {
		return nil, fmt.Errorf("opinion not found: %w", err)
	}

	var supports []models.Opinion
	if err := s.db.Model(&opinion).Association("Supports").Find(&supports); err != nil {
		return nil, fmt.Errorf("failed to get supports: %w", err)
	}

	result := make([]*models.Opinion, len(supports))
	for i := range supports {
		result[i] = &supports[i]
	}
	return result, nil
}

// GetSupportedBy retrieves all opinions that support a given opinion
func (s *OpinionService) GetSupportedBy(opinionID uint) ([]*models.Opinion, error) {
	var opinion models.Opinion
	if err := s.db.First(&opinion, opinionID).Error; err != nil {
		return nil, fmt.Errorf("opinion not found: %w", err)
	}

	var supportedBy []models.Opinion
	if err := s.db.Model(&opinion).Association("SupportedBy").Find(&supportedBy); err != nil {
		return nil, fmt.Errorf("failed to get supported by: %w", err)
	}

	result := make([]*models.Opinion, len(supportedBy))
	for i := range supportedBy {
		result[i] = &supportedBy[i]
	}
	return result, nil
}

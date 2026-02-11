package services

import (
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

type StanceService struct {
	db *gorm.DB
}

func NewStanceService(db *gorm.DB) *StanceService {
	return &StanceService{db: db}
}

// CreateStance creates a new stance (position) on an opinion or reference
func (s *StanceService) CreateStance(userID uint, stanceType, reason string, opinionID, referenceID *uint) (*models.Stance, error) {
	if stanceType != "support" && stanceType != "oppose" && stanceType != "doubt" {
		return nil, fmt.Errorf("invalid stance type: %s (must be 'support', 'oppose', or 'doubt')", stanceType)
	}

	if opinionID == nil && referenceID == nil {
		return nil, fmt.Errorf("stance must target an opinion or a reference")
	}

	stance := &models.Stance{
		UserID:      userID,
		StanceType:  stanceType,
		Reason:      reason,
		OpinionID:   opinionID,
		ReferenceID: referenceID,
	}

	if err := s.db.Create(stance).Error; err != nil {
		return nil, fmt.Errorf("failed to create stance: %w", err)
	}

	return s.GetStanceByID(stance.ID)
}

// GetStanceByID retrieves a stance by ID with user information
func (s *StanceService) GetStanceByID(id uint) (*models.Stance, error) {
	var stance models.Stance
	if err := s.db.Preload("User").First(&stance, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get stance: %w", err)
	}

	stance.Username = stance.User.Username
	return &stance, nil
}

// ListStancesByOpinion retrieves all stances for an opinion
func (s *StanceService) ListStancesByOpinion(opinionID uint) ([]*models.Stance, error) {
	var stances []*models.Stance
	if err := s.db.Preload("User").
		Where("opinion_id = ?", opinionID).
		Order("created_at DESC").
		Find(&stances).Error; err != nil {
		return nil, fmt.Errorf("failed to list stances: %w", err)
	}

	for _, stance := range stances {
		stance.Username = stance.User.Username
	}

	return stances, nil
}

// ListStancesByReference retrieves all stances for a reference
func (s *StanceService) ListStancesByReference(referenceID uint) ([]*models.Stance, error) {
	var stances []*models.Stance
	if err := s.db.Preload("User").
		Where("reference_id = ?", referenceID).
		Order("created_at DESC").
		Find(&stances).Error; err != nil {
		return nil, fmt.Errorf("failed to list stances: %w", err)
	}

	for _, stance := range stances {
		stance.Username = stance.User.Username
	}

	return stances, nil
}

// UpdateStance updates a stance
func (s *StanceService) UpdateStance(id uint, stanceType, reason string) (*models.Stance, error) {
	var stance models.Stance
	if err := s.db.First(&stance, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find stance: %w", err)
	}

	if stanceType != "" {
		if stanceType != "support" && stanceType != "oppose" && stanceType != "doubt" {
			return nil, fmt.Errorf("invalid stance type: %s (must be 'support', 'oppose', or 'doubt')", stanceType)
		}
		stance.StanceType = stanceType
	}
	if reason != "" {
		stance.Reason = reason
	}

	if err := s.db.Save(&stance).Error; err != nil {
		return nil, fmt.Errorf("failed to update stance: %w", err)
	}

	return s.GetStanceByID(id)
}

// DeleteStance soft deletes a stance
func (s *StanceService) DeleteStance(id uint) error {
	if err := s.db.Delete(&models.Stance{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete stance: %w", err)
	}
	return nil
}

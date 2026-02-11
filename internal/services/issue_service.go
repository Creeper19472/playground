package services

import (
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

type IssueService struct {
	db *gorm.DB
}

func NewIssueService(db *gorm.DB) *IssueService {
	return &IssueService{db: db}
}

// CreateIssue creates a new issue
func (s *IssueService) CreateIssue(userID uint, summary, description string) (*models.Issue, error) {
	issue := &models.Issue{
		UserID:      userID,
		Summary:     summary,
		Description: description,
	}
	
	if err := s.db.Create(issue).Error; err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return s.GetIssueByID(issue.ID)
}

// UpdateIssue updates an issue
func (s *IssueService) UpdateIssue(id uint, summary, description string) (*models.Issue, error) {
	var issue models.Issue
	if err := s.db.First(&issue, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find issue: %w", err)
	}

	if summary != "" {
		issue.Summary = summary
	}
	if description != "" {
		issue.Description = description
	}

	if err := s.db.Save(&issue).Error; err != nil {
		return nil, fmt.Errorf("failed to update issue: %w", err)
	}

	return s.GetIssueByID(id)
}

// DeleteIssue soft deletes an issue
func (s *IssueService) DeleteIssue(id uint) error {
	if err := s.db.Delete(&models.Issue{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}
	return nil
}

// GetIssueByID retrieves an issue by ID with user information
func (s *IssueService) GetIssueByID(id uint) (*models.Issue, error) {
	var issue models.Issue
	if err := s.db.Preload("User").First(&issue, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	
	// Set username from loaded User
	issue.Username = issue.User.Username
	
	return &issue, nil
}

// ListIssues retrieves all issues, ordered by vote count
func (s *IssueService) ListIssues() ([]*models.Issue, error) {
	var issues []*models.Issue
	if err := s.db.Preload("User").
		Order("vote_count DESC, created_at DESC").
		Find(&issues).Error; err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}
	
	// Set usernames from loaded Users
	for _, issue := range issues {
		issue.Username = issue.User.Username
	}
	
	return issues, nil
}

// VoteIssue allows a user to vote on an issue
func (s *IssueService) VoteIssue(userID, issueID uint) error {
	// Check if user already voted
	var count int64
	if err := s.db.Model(&models.Vote{}).
		Where("user_id = ? AND issue_id = ?", userID, issueID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing vote: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("user already voted on this issue")
	}

	// Use transaction for atomic operation
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Insert vote
		vote := &models.Vote{
			UserID:  userID,
			IssueID: issueID,
		}
		if err := tx.Create(vote).Error; err != nil {
			return fmt.Errorf("failed to insert vote: %w", err)
		}

		// Update vote count
		if err := tx.Model(&models.Issue{}).
			Where("id = ?", issueID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", 1)).Error; err != nil {
			return fmt.Errorf("failed to update vote count: %w", err)
		}

		return nil
	})
}

// UnvoteIssue allows a user to remove their vote from an issue
func (s *IssueService) UnvoteIssue(userID, issueID uint) error {
	// Check if user has voted
	var count int64
	if err := s.db.Model(&models.Vote{}).
		Where("user_id = ? AND issue_id = ?", userID, issueID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing vote: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("user has not voted on this issue")
	}

	// Use transaction for atomic operation
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete vote
		if err := tx.Where("user_id = ? AND issue_id = ?", userID, issueID).
			Delete(&models.Vote{}).Error; err != nil {
			return fmt.Errorf("failed to delete vote: %w", err)
		}

		// Update vote count
		if err := tx.Model(&models.Issue{}).
			Where("id = ?", issueID).
			UpdateColumn("vote_count", gorm.Expr("vote_count - ?", 1)).Error; err != nil {
			return fmt.Errorf("failed to update vote count: %w", err)
		}

		return nil
	})
}

// GetUserVote checks if a user has voted on an issue
func (s *IssueService) GetUserVote(userID, issueID uint) (bool, error) {
	var count int64
	if err := s.db.Model(&models.Vote{}).
		Where("user_id = ? AND issue_id = ?", userID, issueID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user vote: %w", err)
	}
	return count > 0, nil
}

package services

import (
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

type MessageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(userID uint, content string, issueID *uint) (*models.Message, error) {
	message := &models.Message{
		UserID:  userID,
		Content: content,
		IssueID: issueID,
	}
	
	if err := s.db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return s.GetMessageByID(message.ID)
}

// GetMessageByID retrieves a message by ID with user information
func (s *MessageService) GetMessageByID(id uint) (*models.Message, error) {
	var message models.Message
	if err := s.db.Preload("User").First(&message, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	
	// Set username from loaded User
	message.Username = message.User.Username
	
	return &message, nil
}

// ListMessages retrieves all messages (optionally filtered by issue)
func (s *MessageService) ListMessages(issueID *uint, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	query := s.db.Preload("User").Order("created_at DESC")
	
	if issueID != nil {
		query = query.Where("issue_id = ?", *issueID)
	}
	
	if err := query.Limit(limit).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	
	// Set usernames from loaded Users
	for _, message := range messages {
		message.Username = message.User.Username
	}
	
	return messages, nil
}

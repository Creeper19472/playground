package services

import (
	"database/sql"
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
)

type MessageService struct {
	db *sql.DB
}

func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{db: db}
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(userID int64, content string, issueID *int64) (*models.Message, error) {
	result, err := s.db.Exec(
		"INSERT INTO messages (user_id, content, issue_id) VALUES (?, ?, ?)",
		userID, content, issueID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID: %w", err)
	}

	return s.GetMessageByID(id)
}

// GetMessageByID retrieves a message by ID
func (s *MessageService) GetMessageByID(id int64) (*models.Message, error) {
	message := &models.Message{}
	err := s.db.QueryRow(`
		SELECT m.id, m.user_id, u.username, m.content, m.issue_id, m.created_at
		FROM messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.id = ?`,
		id,
	).Scan(&message.ID, &message.UserID, &message.Username, &message.Content, &message.IssueID, &message.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return message, nil
}

// ListMessages retrieves all messages (optionally filtered by issue)
func (s *MessageService) ListMessages(issueID *int64, limit int) ([]*models.Message, error) {
	var rows *sql.Rows
	var err error

	if issueID != nil {
		rows, err = s.db.Query(`
			SELECT m.id, m.user_id, u.username, m.content, m.issue_id, m.created_at
			FROM messages m
			JOIN users u ON m.user_id = u.id
			WHERE m.issue_id = ?
			ORDER BY m.created_at DESC
			LIMIT ?`,
			issueID, limit,
		)
	} else {
		rows, err = s.db.Query(`
			SELECT m.id, m.user_id, u.username, m.content, m.issue_id, m.created_at
			FROM messages m
			JOIN users u ON m.user_id = u.id
			ORDER BY m.created_at DESC
			LIMIT ?`,
			limit,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		if err := rows.Scan(&message.ID, &message.UserID, &message.Username, &message.Content, &message.IssueID, &message.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

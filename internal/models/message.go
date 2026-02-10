package models

import "time"

// Message represents a chat message
type Message struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	IssueID   *int64    `json:"issue_id,omitempty" db:"issue_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

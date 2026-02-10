package models

import "time"

// Reference represents a reference/link to support an argument
type Reference struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	IssueID     int64     `json:"issue_id" db:"issue_id"`
	URL         string    `json:"url" db:"url"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

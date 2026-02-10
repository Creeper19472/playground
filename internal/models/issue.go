package models

import "time"

// Issue represents a user-proposed issue with voting support
type Issue struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Username    string    `json:"username" db:"username"`
	Summary     string    `json:"summary" db:"summary"`
	Description string    `json:"description" db:"description"`
	VoteCount   int       `json:"vote_count" db:"vote_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Vote represents a user's vote on an issue
type Vote struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	IssueID   int64     `json:"issue_id" db:"issue_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// Issue represents a user-proposed issue with voting support
type Issue struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	Username    string         `json:"username" gorm:"-"` // Computed field
	Summary     string         `json:"summary" gorm:"not null"`
	Description string         `json:"description"`
	VoteCount   int            `json:"vote_count" gorm:"default:0;index:idx_vote_count,sort:desc"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Votes      []Vote      `json:"-" gorm:"foreignKey:IssueID"`
	Messages   []Message   `json:"-" gorm:"foreignKey:IssueID"`
	References []Reference `json:"-" gorm:"foreignKey:IssueID"`
}

// Vote represents a user's vote on an issue
type Vote struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_user_issue"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	IssueID   uint           `json:"issue_id" gorm:"not null;uniqueIndex:idx_user_issue;index"`
	Issue     Issue          `json:"-" gorm:"foreignKey:IssueID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

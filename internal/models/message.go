package models

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a chat message
type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Username  string         `json:"username" gorm:"-"` // Computed field
	Content   string         `json:"content" gorm:"not null"`
	IssueID   *uint          `json:"issue_id,omitempty" gorm:"index"`
	Issue     *Issue         `json:"-" gorm:"foreignKey:IssueID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

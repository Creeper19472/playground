package models

import (
	"time"

	"gorm.io/gorm"
)

// Reference represents a reference/link to support an argument
type Reference struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	IssueID     uint           `json:"issue_id" gorm:"not null;index"`
	Issue       Issue          `json:"-" gorm:"foreignKey:IssueID"`
	URL         string         `json:"url" gorm:"not null"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

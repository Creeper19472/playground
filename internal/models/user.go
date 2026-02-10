package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Issues    []Issue    `json:"-" gorm:"foreignKey:UserID"`
	Messages  []Message  `json:"-" gorm:"foreignKey:UserID"`
	Votes     []Vote     `json:"-" gorm:"foreignKey:UserID"`
	References []Reference `json:"-" gorm:"foreignKey:UserID"`
}

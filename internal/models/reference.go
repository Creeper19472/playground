package models

import (
	"time"

	"gorm.io/gorm"
)

// Reference represents a fact/reference source to support an opinion.
// A reference can be linked to an issue directly and/or to a specific opinion.
type Reference struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	IssueID     uint           `json:"issue_id" gorm:"not null;index"`
	Issue       Issue          `json:"-" gorm:"foreignKey:IssueID"`
	OpinionID   *uint          `json:"opinion_id,omitempty" gorm:"index"` // Optional link to an opinion
	Opinion     *Opinion       `json:"-" gorm:"foreignKey:OpinionID"`
	URL         string         `json:"url" gorm:"not null"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Stances on this reference
	Stances []Stance `json:"-" gorm:"foreignKey:ReferenceID"`
}

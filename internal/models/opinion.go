package models

import (
	"time"

	"gorm.io/gorm"
)

// Opinion represents a user's opinion under an issue.
// Opinions can have hierarchical relationships: one opinion can support another opinion,
// multiple opinions can support one opinion, and one opinion can support multiple opinions.
type Opinion struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Username  string         `json:"username" gorm:"-"` // Computed field
	IssueID   uint           `json:"issue_id" gorm:"not null;index"`
	Issue     Issue          `json:"-" gorm:"foreignKey:IssueID"`
	Content   string         `json:"content" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Many-to-many self-referencing: opinions that this opinion supports
	Supports []Opinion `json:"-" gorm:"many2many:opinion_supports;joinForeignKey:SupporterID;joinReferences:SupportedID"`
	// Many-to-many self-referencing: opinions that support this opinion
	SupportedBy []Opinion `json:"-" gorm:"many2many:opinion_supports;joinForeignKey:SupportedID;joinReferences:SupporterID"`

	// References (facts) attached to this opinion
	References []Reference `json:"-" gorm:"foreignKey:OpinionID"`
	// Stances on this opinion
	Stances []Stance `json:"-" gorm:"foreignKey:OpinionID"`
}

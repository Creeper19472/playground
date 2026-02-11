package models

import (
	"time"

	"gorm.io/gorm"
)

// Issue represents a user-proposed issue or proposition.
// When Type is "issue", users determine a conclusion through discussion and voting.
// When Type is "proposition", it has a definite truth value (true or false).
type Issue struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	Username    string         `json:"username" gorm:"-"` // Computed field
	Type        string         `json:"type" gorm:"not null;default:'issue';index"` // "issue" or "proposition"
	Summary     string         `json:"summary" gorm:"not null"`
	Description string         `json:"description"`
	Conclusion  string         `json:"conclusion,omitempty"` // Final conclusion for issues (determined by discussion/voting)
	TruthValue  *bool          `json:"truth_value,omitempty"` // For propositions: true or false
	VoteCount   int            `json:"vote_count" gorm:"default:0;index:idx_vote_count,sort:desc"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Votes      []Vote      `json:"-" gorm:"foreignKey:IssueID"`
	Messages   []Message   `json:"-" gorm:"foreignKey:IssueID"`
	References []Reference `json:"-" gorm:"foreignKey:IssueID"`
	Opinions   []Opinion   `json:"-" gorm:"foreignKey:IssueID"`
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

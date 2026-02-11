package models

import (
	"time"

	"gorm.io/gorm"
)

// Stance represents a user's position on an opinion or a fact (reference).
// Users can express their doubts or agreement by "stating their position"
// and can attach their reasons for believing so.
type Stance struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	Username    string         `json:"username" gorm:"-"` // Computed field
	StanceType  string         `json:"stance_type" gorm:"not null"` // "support", "oppose", "doubt"
	Reason      string         `json:"reason"` // Reason for taking this stance
	OpinionID   *uint          `json:"opinion_id,omitempty" gorm:"index"` // Stance on an opinion
	Opinion     *Opinion       `json:"-" gorm:"foreignKey:OpinionID"`
	ReferenceID *uint          `json:"reference_id,omitempty" gorm:"index"` // Stance on a fact/reference
	Reference   *Reference     `json:"-" gorm:"foreignKey:ReferenceID"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

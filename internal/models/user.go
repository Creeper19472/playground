package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"not null"` // Never expose in JSON
	IsAdmin      bool           `json:"is_admin" gorm:"default:false"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Issues      []Issue      `json:"-" gorm:"foreignKey:UserID"`
	Messages    []Message    `json:"-" gorm:"foreignKey:UserID"`
	Votes       []Vote       `json:"-" gorm:"foreignKey:UserID"`
	References  []Reference  `json:"-" gorm:"foreignKey:UserID"`
	Opinions    []Opinion    `json:"-" gorm:"foreignKey:UserID"`
	Stances     []Stance     `json:"-" gorm:"foreignKey:UserID"`
	Roles       []Role       `json:"-" gorm:"many2many:user_roles;"`
	Groups      []UserGroup  `json:"-" gorm:"many2many:user_group_members;"`
	Permissions []Permission `json:"-" gorm:"many2many:user_permissions;"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies if the provided password matches the hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a user role (e.g., admin, moderator, user)
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Users       []User       `json:"-" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"-" gorm:"many2many:role_permissions;"`
}

// Permission represents a specific permission (e.g., create_issue, delete_user)
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Resource    string         `json:"resource" gorm:"not null;index"` // e.g., "issue", "user", "message"
	Action      string         `json:"action" gorm:"not null;index"`   // e.g., "create", "read", "update", "delete"
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Roles  []Role      `json:"-" gorm:"many2many:role_permissions;"`
	Users  []User      `json:"-" gorm:"many2many:user_permissions;"`
	Groups []UserGroup `json:"-" gorm:"many2many:group_permissions;"`
}

// UserGroup represents a group of users that can inherit permissions
type UserGroup struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Users       []User       `json:"-" gorm:"many2many:user_group_members;"`
	Permissions []Permission `json:"-" gorm:"many2many:group_permissions;"`
}

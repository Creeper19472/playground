package services

import (
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user (for backwards compatibility, without password)
func (s *UserService) CreateUser(username, email string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Email:    email,
		IsActive: true,
	}
	
	// Set a default password that must be changed
	if err := user.SetPassword("changeme123"); err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}
	
	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() ([]*models.User, error) {
	var users []*models.User
	if err := s.db.Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id uint, username, email string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// SetUserActiveStatus sets whether a user is active
func (s *UserService) SetUserActiveStatus(id uint, isActive bool) error {
	if err := s.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// SetUserAdminStatus sets whether a user is an admin
func (s *UserService) SetUserAdminStatus(id uint, isAdmin bool) error {
	if err := s.db.Model(&models.User{}).Where("id = ?", id).Update("is_admin", isAdmin).Error; err != nil {
		return fmt.Errorf("failed to update admin status: %w", err)
	}
	return nil
}


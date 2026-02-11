package services

import (
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/gorm"
)

type PermissionService struct {
	db *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// HasPermission checks if a user has a specific permission (directly, via role, or via group)
func (s *PermissionService) HasPermission(userID uint, resource, action string) (bool, error) {
	var user models.User
	if err := s.db.Preload("Permissions").Preload("Roles.Permissions").Preload("Groups.Permissions").First(&user, userID).Error; err != nil {
		return false, fmt.Errorf("failed to load user: %w", err)
	}

	// Admin users have all permissions
	if user.IsAdmin {
		return true, nil
	}

	// Check direct user permissions
	for _, perm := range user.Permissions {
		if perm.Resource == resource && perm.Action == action {
			return true, nil
		}
	}

	// Check role permissions
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm.Resource == resource && perm.Action == action {
				return true, nil
			}
		}
	}

	// Check group permissions
	for _, group := range user.Groups {
		for _, perm := range group.Permissions {
			if perm.Resource == resource && perm.Action == action {
				return true, nil
			}
		}
	}

	return false, nil
}

// GrantPermissionToUser grants a permission directly to a user
func (s *PermissionService) GrantPermissionToUser(userID, permissionID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	return s.db.Model(&user).Association("Permissions").Append(&permission)
}

// RevokePermissionFromUser revokes a permission from a user
func (s *PermissionService) RevokePermissionFromUser(userID, permissionID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	return s.db.Model(&user).Association("Permissions").Delete(&permission)
}

// AssignRoleToUser assigns a role to a user
func (s *PermissionService) AssignRoleToUser(userID, roleID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return s.db.Model(&user).Association("Roles").Append(&role)
}

// RemoveRoleFromUser removes a role from a user
func (s *PermissionService) RemoveRoleFromUser(userID, roleID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return s.db.Model(&user).Association("Roles").Delete(&role)
}

// AddUserToGroup adds a user to a group
func (s *PermissionService) AddUserToGroup(userID, groupID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var group models.UserGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	return s.db.Model(&user).Association("Groups").Append(&group)
}

// RemoveUserFromGroup removes a user from a group
func (s *PermissionService) RemoveUserFromGroup(userID, groupID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var group models.UserGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	return s.db.Model(&user).Association("Groups").Delete(&group)
}

// CreatePermission creates a new permission
func (s *PermissionService) CreatePermission(name, description, resource, action string) (*models.Permission, error) {
	permission := &models.Permission{
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
	}

	if err := s.db.Create(permission).Error; err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

// CreateRole creates a new role
func (s *PermissionService) CreateRole(name, description string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// CreateUserGroup creates a new user group
func (s *PermissionService) CreateUserGroup(name, description string) (*models.UserGroup, error) {
	group := &models.UserGroup{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(group).Error; err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

// GrantPermissionToRole grants a permission to a role
func (s *PermissionService) GrantPermissionToRole(roleID, permissionID uint) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	return s.db.Model(&role).Association("Permissions").Append(&permission)
}

// GrantPermissionToGroup grants a permission to a group
func (s *PermissionService) GrantPermissionToGroup(groupID, permissionID uint) error {
	var group models.UserGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	return s.db.Model(&group).Association("Permissions").Append(&permission)
}

// ListPermissions returns all permissions
func (s *PermissionService) ListPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	return permissions, nil
}

// ListRoles returns all roles
func (s *PermissionService) ListRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, nil
}

// ListUserGroups returns all user groups
func (s *PermissionService) ListUserGroups() ([]models.UserGroup, error) {
	var groups []models.UserGroup
	if err := s.db.Preload("Permissions").Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	return groups, nil
}

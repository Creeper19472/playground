package services

import (
	"database/sql"
	"fmt"

	"github.com/Creeper19472/playground/internal/models"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, email string) (*models.User, error) {
	result, err := s.db.Exec(
		"INSERT INTO users (username, email) VALUES (?, ?)",
		username, email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	return s.GetUserByID(id)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(
		"SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(
		"SELECT id, username, email, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() ([]*models.User, error) {
	rows, err := s.db.Query(
		"SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

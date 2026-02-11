package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Creeper19472/playground/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user is not active")
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret []byte
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}
	return &AuthService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	var user models.User
	
	// Find user by username or email
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return "", nil, ErrUserNotActive
	}

	// Verify password
	if !user.CheckPassword(password) {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, &user, nil
}

// GenerateToken creates a new JWT token for a user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "playground",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Register creates a new user with password
func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Email:    email,
		IsActive: true,
		IsAdmin:  false,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if !user.CheckPassword(oldPassword) {
		return ErrInvalidCredentials
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

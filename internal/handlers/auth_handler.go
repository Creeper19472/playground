package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Creeper19472/playground/internal/middleware"
	"github.com/Creeper19472/playground/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		case services.ErrUserNotActive:
			http.Error(w, "User is not active", http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password == "" || len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		http.Error(w, "New password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			http.Error(w, "Invalid old password", http.StatusUnauthorized)
		case services.ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a JWT-based system, logout is handled client-side by deleting the token
	// Here we just return a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	username, _ := middleware.GetUsername(r)
	isAdmin := middleware.IsAdmin(r)

	response := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"is_admin": isAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

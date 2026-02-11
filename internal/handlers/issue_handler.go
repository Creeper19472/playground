package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Creeper19472/playground/internal/middleware"
	"github.com/Creeper19472/playground/internal/services"
	ws "github.com/Creeper19472/playground/pkg/websocket"
	"github.com/gorilla/mux"
)

type IssueHandler struct {
	issueService *services.IssueService
	hub          *ws.Hub
}

func NewIssueHandler(issueService *services.IssueService, hub *ws.Hub) *IssueHandler {
	return &IssueHandler{
		issueService: issueService,
		hub:          hub,
	}
}

type CreateIssueRequest struct {
	UserID      uint   `json:"user_id"`
	Type        string `json:"type"` // "issue" or "proposition"
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type UpdateIssueRequest struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Conclusion  string `json:"conclusion,omitempty"`
	TruthValue  *bool  `json:"truth_value,omitempty"`
}

type VoteRequest struct {
	UserID uint `json:"user_id"`
}

func (h *IssueHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var req CreateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	issue, err := h.issueService.CreateIssue(req.UserID, req.Type, req.Summary, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast to all connected clients
	h.hub.BroadcastMessage("new_issue", issue)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func (h *IssueHandler) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	// Get the issue to check ownership
	issue, err := h.issueService.GetIssueByID(uint(id))
	if err != nil {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	// Check if user can update (must be owner or admin)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if issue.UserID != userID && !middleware.IsAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedIssue, err := h.issueService.UpdateIssue(uint(id), req.Summary, req.Description, req.Conclusion, req.TruthValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast update
	h.hub.BroadcastMessage("issue_updated", updatedIssue)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedIssue)
}

func (h *IssueHandler) DeleteIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	// Get the issue to check ownership
	issue, err := h.issueService.GetIssueByID(uint(id))
	if err != nil {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	// Check if user can delete (must be owner or admin)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if issue.UserID != userID && !middleware.IsAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.issueService.DeleteIssue(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast deletion
	h.hub.BroadcastMessage("issue_deleted", map[string]interface{}{"id": id})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Issue deleted successfully"})
}

func (h *IssueHandler) GetIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	issue, err := h.issueService.GetIssueByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func (h *IssueHandler) ListIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := h.issueService.ListIssues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issues)
}

func (h *IssueHandler) VoteIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issueID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.issueService.VoteIssue(req.UserID, uint(issueID)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated issue and broadcast
	issue, err := h.issueService.GetIssueByID(uint(issueID))
	if err != nil {
		log.Printf("Failed to get issue after voting: %v", err)
	} else {
		h.hub.BroadcastMessage("issue_voted", issue)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "voted"})
}

func (h *IssueHandler) UnvoteIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issueID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.issueService.UnvoteIssue(req.UserID, uint(issueID)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated issue and broadcast
	issue, err := h.issueService.GetIssueByID(uint(issueID))
	if err != nil {
		log.Printf("Failed to get issue after unvoting: %v", err)
	} else {
		h.hub.BroadcastMessage("issue_unvoted", issue)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "unvoted"})
}

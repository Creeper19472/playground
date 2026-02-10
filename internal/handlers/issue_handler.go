package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	UserID      uint  `json:"user_id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
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

	issue, err := h.issueService.CreateIssue(req.UserID, req.Summary, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast to all connected clients
	h.hub.BroadcastMessage("new_issue", issue)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
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

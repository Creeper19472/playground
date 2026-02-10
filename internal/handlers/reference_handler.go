package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Creeper19472/playground/internal/services"
	"github.com/gorilla/mux"
)

type ReferenceHandler struct {
	referenceService *services.ReferenceService
}

func NewReferenceHandler(referenceService *services.ReferenceService) *ReferenceHandler {
	return &ReferenceHandler{referenceService: referenceService}
}

type CreateReferenceRequest struct {
	UserID      int64  `json:"user_id"`
	IssueID     int64  `json:"issue_id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *ReferenceHandler) CreateReference(w http.ResponseWriter, r *http.Request) {
	var req CreateReferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reference, err := h.referenceService.CreateReference(req.UserID, req.IssueID, req.URL, req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reference)
}

func (h *ReferenceHandler) GetReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid reference ID", http.StatusBadRequest)
		return
	}

	reference, err := h.referenceService.GetReferenceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reference)
}

func (h *ReferenceHandler) ListReferencesByIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issueID, err := strconv.ParseInt(vars["issue_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	references, err := h.referenceService.ListReferencesByIssue(issueID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(references)
}

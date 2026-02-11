package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Creeper19472/playground/internal/middleware"
	"github.com/Creeper19472/playground/internal/services"
	ws "github.com/Creeper19472/playground/pkg/websocket"
	"github.com/gorilla/mux"
)

type OpinionHandler struct {
	opinionService *services.OpinionService
	hub            *ws.Hub
}

func NewOpinionHandler(opinionService *services.OpinionService, hub *ws.Hub) *OpinionHandler {
	return &OpinionHandler{
		opinionService: opinionService,
		hub:            hub,
	}
}

type CreateOpinionRequest struct {
	UserID  uint   `json:"user_id"`
	IssueID uint   `json:"issue_id"`
	Content string `json:"content"`
}

type UpdateOpinionRequest struct {
	Content string `json:"content"`
}

type OpinionSupportRequest struct {
	SupporterID uint `json:"supporter_id"`
	SupportedID uint `json:"supported_id"`
}

func (h *OpinionHandler) CreateOpinion(w http.ResponseWriter, r *http.Request) {
	var req CreateOpinionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opinion, err := h.opinionService.CreateOpinion(req.UserID, req.IssueID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.hub.BroadcastMessage("new_opinion", opinion)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opinion)
}

func (h *OpinionHandler) GetOpinion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid opinion ID", http.StatusBadRequest)
		return
	}

	opinion, err := h.opinionService.GetOpinionByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opinion)
}

func (h *OpinionHandler) ListOpinionsByIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issueID, err := strconv.ParseUint(vars["issue_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	opinions, err := h.opinionService.ListOpinionsByIssue(uint(issueID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opinions)
}

func (h *OpinionHandler) UpdateOpinion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid opinion ID", http.StatusBadRequest)
		return
	}

	// Get the opinion to check ownership
	opinion, err := h.opinionService.GetOpinionByID(uint(id))
	if err != nil {
		http.Error(w, "Opinion not found", http.StatusNotFound)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if opinion.UserID != userID && !middleware.IsAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateOpinionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedOpinion, err := h.opinionService.UpdateOpinion(uint(id), req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.hub.BroadcastMessage("opinion_updated", updatedOpinion)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOpinion)
}

func (h *OpinionHandler) DeleteOpinion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid opinion ID", http.StatusBadRequest)
		return
	}

	opinion, err := h.opinionService.GetOpinionByID(uint(id))
	if err != nil {
		http.Error(w, "Opinion not found", http.StatusNotFound)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if opinion.UserID != userID && !middleware.IsAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.opinionService.DeleteOpinion(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.hub.BroadcastMessage("opinion_deleted", map[string]interface{}{"id": id})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Opinion deleted successfully"})
}

func (h *OpinionHandler) AddSupport(w http.ResponseWriter, r *http.Request) {
	var req OpinionSupportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.opinionService.AddSupport(req.SupporterID, req.SupportedID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Support relationship added"})
}

func (h *OpinionHandler) RemoveSupport(w http.ResponseWriter, r *http.Request) {
	var req OpinionSupportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.opinionService.RemoveSupport(req.SupporterID, req.SupportedID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Support relationship removed"})
}

func (h *OpinionHandler) GetSupports(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid opinion ID", http.StatusBadRequest)
		return
	}

	supports, err := h.opinionService.GetSupports(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supports)
}

func (h *OpinionHandler) GetSupportedBy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid opinion ID", http.StatusBadRequest)
		return
	}

	supportedBy, err := h.opinionService.GetSupportedBy(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supportedBy)
}

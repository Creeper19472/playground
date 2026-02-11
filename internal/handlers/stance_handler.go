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

type StanceHandler struct {
	stanceService *services.StanceService
	hub           *ws.Hub
}

func NewStanceHandler(stanceService *services.StanceService, hub *ws.Hub) *StanceHandler {
	return &StanceHandler{
		stanceService: stanceService,
		hub:           hub,
	}
}

type CreateStanceRequest struct {
	UserID      uint   `json:"user_id"`
	StanceType  string `json:"stance_type"` // "support", "oppose", "doubt"
	Reason      string `json:"reason"`
	OpinionID   *uint  `json:"opinion_id,omitempty"`
	ReferenceID *uint  `json:"reference_id,omitempty"`
}

type UpdateStanceRequest struct {
	StanceType string `json:"stance_type,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

func (h *StanceHandler) CreateStance(w http.ResponseWriter, r *http.Request) {
	var req CreateStanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stance, err := h.stanceService.CreateStance(req.UserID, req.StanceType, req.Reason, req.OpinionID, req.ReferenceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.hub.BroadcastMessage("new_stance", stance)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stance)
}

func (h *StanceHandler) GetStance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid stance ID", http.StatusBadRequest)
		return
	}

	stance, err := h.stanceService.GetStanceByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stance)
}

func (h *StanceHandler) ListStancesByOpinion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	opinionID, err := strconv.ParseUint(vars["opinion_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid opinion ID", http.StatusBadRequest)
		return
	}

	stances, err := h.stanceService.ListStancesByOpinion(uint(opinionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stances)
}

func (h *StanceHandler) ListStancesByReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	referenceID, err := strconv.ParseUint(vars["reference_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid reference ID", http.StatusBadRequest)
		return
	}

	stances, err := h.stanceService.ListStancesByReference(uint(referenceID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stances)
}

func (h *StanceHandler) UpdateStance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid stance ID", http.StatusBadRequest)
		return
	}

	stance, err := h.stanceService.GetStanceByID(uint(id))
	if err != nil {
		http.Error(w, "Stance not found", http.StatusNotFound)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if stance.UserID != userID && !middleware.IsAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateStanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedStance, err := h.stanceService.UpdateStance(uint(id), req.StanceType, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.hub.BroadcastMessage("stance_updated", updatedStance)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStance)
}

func (h *StanceHandler) DeleteStance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid stance ID", http.StatusBadRequest)
		return
	}

	stance, err := h.stanceService.GetStanceByID(uint(id))
	if err != nil {
		http.Error(w, "Stance not found", http.StatusNotFound)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if stance.UserID != userID && !middleware.IsAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.stanceService.DeleteStance(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.hub.BroadcastMessage("stance_deleted", map[string]interface{}{"id": id})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Stance deleted successfully"})
}

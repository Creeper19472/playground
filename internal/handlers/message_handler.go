package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Creeper19472/playground/internal/services"
	ws "github.com/Creeper19472/playground/pkg/websocket"
	"github.com/gorilla/mux"
)

type MessageHandler struct {
	messageService *services.MessageService
	hub            *ws.Hub
}

func NewMessageHandler(messageService *services.MessageService, hub *ws.Hub) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		hub:            hub,
	}
}

type CreateMessageRequest struct {
	UserID  uint  `json:"user_id"`
	Content string `json:"content"`
	IssueID *uint `json:"issue_id,omitempty"`
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := h.messageService.CreateMessage(req.UserID, req.Content, req.IssueID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast to all connected clients
	h.hub.BroadcastMessage("new_message", message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	message, err := h.messageService.GetMessageByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *MessageHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	issueIDStr := r.URL.Query().Get("issue_id")
	limitStr := r.URL.Query().Get("limit")

	var issueID *uint
	if issueIDStr != "" {
		id, err := strconv.ParseUint(issueIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid issue ID", http.StatusBadRequest)
			return
		}
		uid := uint(id)
		issueID = &uid
	}

	limit := 100
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = l
	}

	messages, err := h.messageService.ListMessages(issueID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

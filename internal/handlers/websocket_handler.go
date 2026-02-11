package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	ws "github.com/Creeper19472/playground/pkg/websocket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get user info from query parameters
	userIDStr := r.URL.Query().Get("user_id")
	username := r.URL.Query().Get("username")

	if userIDStr == "" || username == "" {
		http.Error(w, "user_id and username required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &ws.Client{
		ID:       uuid.New().String(),
		UserID:   uint(userID),
		Username: username,
		Conn:     conn,
		Hub:      h.hub,
		Send:     make(chan []byte, 256),
	}

	h.hub.Register <- client

	// Send welcome message
	welcomeMsg := fmt.Sprintf(`{"type":"connected","data":{"message":"Welcome %s! Connected clients: %d"}}`, username, h.hub.GetClientCount())
	client.Send <- []byte(welcomeMsg)

	// Start read and write pumps
	go client.WritePump()
	go client.ReadPump()
}

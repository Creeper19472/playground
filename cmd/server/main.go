package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Creeper19472/playground/config"
	"github.com/Creeper19472/playground/internal/db"
	"github.com/Creeper19472/playground/internal/handlers"
	"github.com/Creeper19472/playground/internal/services"
	"github.com/Creeper19472/playground/pkg/websocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "", "Path to configuration file (YAML)")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Printf("Loaded configuration: Server=%s:%d, Database=%s", 
		cfg.Server.Host, cfg.Server.Port, cfg.Database.Type)

	// Initialize database
	database, err := db.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize services
	userService := services.NewUserService(database.DB)
	messageService := services.NewMessageService(database.DB)
	issueService := services.NewIssueService(database.DB)
	referenceService := services.NewReferenceService(database.DB)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	messageHandler := handlers.NewMessageHandler(messageService, hub)
	issueHandler := handlers.NewIssueHandler(issueService, hub)
	referenceHandler := handlers.NewReferenceHandler(referenceService)
	wsHandler := handlers.NewWebSocketHandler(hub)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users", userHandler.ListUsers).Methods("GET")

	// Message routes
	api.HandleFunc("/messages", messageHandler.CreateMessage).Methods("POST")
	api.HandleFunc("/messages/{id}", messageHandler.GetMessage).Methods("GET")
	api.HandleFunc("/messages", messageHandler.ListMessages).Methods("GET")

	// Issue routes
	api.HandleFunc("/issues", issueHandler.CreateIssue).Methods("POST")
	api.HandleFunc("/issues/{id}", issueHandler.GetIssue).Methods("GET")
	api.HandleFunc("/issues", issueHandler.ListIssues).Methods("GET")
	api.HandleFunc("/issues/{id}/vote", issueHandler.VoteIssue).Methods("POST")
	api.HandleFunc("/issues/{id}/unvote", issueHandler.UnvoteIssue).Methods("POST")

	// Reference routes
	api.HandleFunc("/references", referenceHandler.CreateReference).Methods("POST")
	api.HandleFunc("/references/{id}", referenceHandler.GetReference).Methods("GET")
	api.HandleFunc("/issues/{issue_id}/references", referenceHandler.ListReferencesByIssue).Methods("GET")

	// WebSocket route
	router.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","connected_clients":%d}`, hub.GetClientCount())
	}).Methods("GET")

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server starting on %s...", addr)
	
	// Log user-friendly URLs
	displayHost := cfg.Server.Host
	if displayHost == "0.0.0.0" || displayHost == "" {
		displayHost = "localhost"
	}
	log.Printf("WebSocket endpoint: ws://%s:%d/ws", displayHost, cfg.Server.Port)
	log.Printf("API endpoint: http://%s:%d/api/v1", displayHost, cfg.Server.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}

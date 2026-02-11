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
	"github.com/Creeper19472/playground/internal/middleware"
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

	// Initialize admin account if it doesn't exist
	adminInit := services.NewAdminInitializer(database.DB)
	created, username, password, err := adminInit.InitializeAdminAccount()
	if err != nil {
		log.Fatalf("Failed to initialize admin account: %v", err)
	}
	
	if created {
		log.Printf("Administrator account created successfully")
		
		// Write credentials to a secure file
		credentialsFile := "admin_credentials.txt"
		if err := services.WriteCredentialsToFile(username, password, credentialsFile); err != nil {
			log.Fatalf("Failed to write admin credentials to file: %v", err)
		}
		
		log.Printf("⚠️  IMPORTANT: Administrator credentials have been written to '%s'", credentialsFile)
		log.Printf("⚠️  Please save these credentials and delete the file for security!")
		log.Printf("⚠️  Change the admin password immediately after first login.")
	} else {
		log.Printf("Administrator account already exists, skipping creation")
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize services
	userService := services.NewUserService(database.DB)
	messageService := services.NewMessageService(database.DB)
	issueService := services.NewIssueService(database.DB)
	referenceService := services.NewReferenceService(database.DB)
	opinionService := services.NewOpinionService(database.DB)
	stanceService := services.NewStanceService(database.DB)
	authService := services.NewAuthService(database.DB, cfg.Auth.JWTSecret)
	permissionService := services.NewPermissionService(database.DB)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	messageHandler := handlers.NewMessageHandler(messageService, hub)
	issueHandler := handlers.NewIssueHandler(issueService, hub)
	referenceHandler := handlers.NewReferenceHandler(referenceService)
	opinionHandler := handlers.NewOpinionHandler(opinionService, hub)
	stanceHandler := handlers.NewStanceHandler(stanceService, hub)
	wsHandler := handlers.NewWebSocketHandler(hub)
	authHandler := handlers.NewAuthHandler(authService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public auth routes (no authentication required)
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Protected auth routes (authentication required)
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.Use(middleware.AuthMiddleware(authService))
	authRoutes.HandleFunc("/me", authHandler.Me).Methods("GET")
	authRoutes.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	authRoutes.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")

	// User routes (protected)
	protectedUserRoutes := api.PathPrefix("/users").Subrouter()
	protectedUserRoutes.Use(middleware.AuthMiddleware(authService))
	protectedUserRoutes.HandleFunc("", userHandler.ListUsers).Methods("GET")
	protectedUserRoutes.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	protectedUserRoutes.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	protectedUserRoutes.Handle("/{id}", middleware.AdminOnlyMiddleware(http.HandlerFunc(userHandler.DeleteUser))).Methods("DELETE")

	// Message routes (authentication required)
	protectedMessageRoutes := api.PathPrefix("/messages").Subrouter()
	protectedMessageRoutes.Use(middleware.AuthMiddleware(authService))
	protectedMessageRoutes.HandleFunc("", messageHandler.CreateMessage).Methods("POST")
	protectedMessageRoutes.HandleFunc("/{id}", messageHandler.GetMessage).Methods("GET")
	protectedMessageRoutes.HandleFunc("", messageHandler.ListMessages).Methods("GET")

	// Issue routes (authentication required for create/vote, public for read)
	api.HandleFunc("/issues", issueHandler.ListIssues).Methods("GET")
	api.HandleFunc("/issues/{id}", issueHandler.GetIssue).Methods("GET")
	
	protectedIssueRoutes := api.PathPrefix("/issues").Subrouter()
	protectedIssueRoutes.Use(middleware.AuthMiddleware(authService))
	protectedIssueRoutes.HandleFunc("", issueHandler.CreateIssue).Methods("POST")
	protectedIssueRoutes.HandleFunc("/{id}", issueHandler.UpdateIssue).Methods("PUT")
	protectedIssueRoutes.HandleFunc("/{id}", issueHandler.DeleteIssue).Methods("DELETE")
	protectedIssueRoutes.HandleFunc("/{id}/vote", issueHandler.VoteIssue).Methods("POST")
	protectedIssueRoutes.HandleFunc("/{id}/unvote", issueHandler.UnvoteIssue).Methods("POST")

	// Reference routes (authentication required)
	protectedReferenceRoutes := api.PathPrefix("/references").Subrouter()
	protectedReferenceRoutes.Use(middleware.AuthMiddleware(authService))
	protectedReferenceRoutes.HandleFunc("", referenceHandler.CreateReference).Methods("POST")
	protectedReferenceRoutes.HandleFunc("/{id}", referenceHandler.GetReference).Methods("GET")
	
	api.HandleFunc("/issues/{issue_id}/references", referenceHandler.ListReferencesByIssue).Methods("GET")

	// Opinion routes (authentication required for create/update/delete, public for read)
	api.HandleFunc("/issues/{issue_id}/opinions", opinionHandler.ListOpinionsByIssue).Methods("GET")
	api.HandleFunc("/opinions/{id}", opinionHandler.GetOpinion).Methods("GET")
	api.HandleFunc("/opinions/{id}/supports", opinionHandler.GetSupports).Methods("GET")
	api.HandleFunc("/opinions/{id}/supported-by", opinionHandler.GetSupportedBy).Methods("GET")

	protectedOpinionRoutes := api.PathPrefix("/opinions").Subrouter()
	protectedOpinionRoutes.Use(middleware.AuthMiddleware(authService))
	protectedOpinionRoutes.HandleFunc("", opinionHandler.CreateOpinion).Methods("POST")
	protectedOpinionRoutes.HandleFunc("/{id}", opinionHandler.UpdateOpinion).Methods("PUT")
	protectedOpinionRoutes.HandleFunc("/{id}", opinionHandler.DeleteOpinion).Methods("DELETE")
	protectedOpinionRoutes.HandleFunc("/support", opinionHandler.AddSupport).Methods("POST")
	protectedOpinionRoutes.HandleFunc("/support", opinionHandler.RemoveSupport).Methods("DELETE")

	// Stance routes (authentication required for create/update/delete, public for read)
	api.HandleFunc("/opinions/{opinion_id}/stances", stanceHandler.ListStancesByOpinion).Methods("GET")
	api.HandleFunc("/references/{reference_id}/stances", stanceHandler.ListStancesByReference).Methods("GET")
	api.HandleFunc("/stances/{id}", stanceHandler.GetStance).Methods("GET")

	protectedStanceRoutes := api.PathPrefix("/stances").Subrouter()
	protectedStanceRoutes.Use(middleware.AuthMiddleware(authService))
	protectedStanceRoutes.HandleFunc("", stanceHandler.CreateStance).Methods("POST")
	protectedStanceRoutes.HandleFunc("/{id}", stanceHandler.UpdateStance).Methods("PUT")
	protectedStanceRoutes.HandleFunc("/{id}", stanceHandler.DeleteStance).Methods("DELETE")

	// Permission management routes (admin only)
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.AuthMiddleware(authService))
	adminRoutes.Use(middleware.AdminOnlyMiddleware)
	
	// Permission CRUD
	adminRoutes.HandleFunc("/permissions", permissionHandler.ListPermissions).Methods("GET")
	adminRoutes.HandleFunc("/permissions", permissionHandler.CreatePermission).Methods("POST")
	adminRoutes.HandleFunc("/roles", permissionHandler.ListRoles).Methods("GET")
	adminRoutes.HandleFunc("/roles", permissionHandler.CreateRole).Methods("POST")
	adminRoutes.HandleFunc("/groups", permissionHandler.ListGroups).Methods("GET")
	adminRoutes.HandleFunc("/groups", permissionHandler.CreateGroup).Methods("POST")
	
	// User permission management
	adminRoutes.HandleFunc("/users/{user_id}/permissions", permissionHandler.GrantPermissionToUser).Methods("POST")
	adminRoutes.HandleFunc("/users/{user_id}/permissions", permissionHandler.RevokePermissionFromUser).Methods("DELETE")
	adminRoutes.HandleFunc("/users/{user_id}/roles", permissionHandler.AssignRoleToUser).Methods("POST")
	adminRoutes.HandleFunc("/users/{user_id}/roles", permissionHandler.RemoveRoleFromUser).Methods("DELETE")
	adminRoutes.HandleFunc("/users/{user_id}/groups", permissionHandler.AddUserToGroup).Methods("POST")
	adminRoutes.HandleFunc("/users/{user_id}/groups", permissionHandler.RemoveUserFromGroup).Methods("DELETE")
	
	// Role and group permission management
	adminRoutes.HandleFunc("/roles/{role_id}/permissions", permissionHandler.GrantPermissionToRole).Methods("POST")
	adminRoutes.HandleFunc("/groups/{group_id}/permissions", permissionHandler.GrantPermissionToGroup).Methods("POST")

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

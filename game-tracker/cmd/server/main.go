package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"game-tracker/internal/api"
	"game-tracker/internal/config"
	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
	"game-tracker/internal/middleware"
	"game-tracker/internal/worker"
)

// frontendFS will be populated by frontend_embed.go when built with embedded frontend
var frontendFS fs.FS

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initialize Firebase Admin SDK
	var opts []option.ClientOption
	if cfg.Firebase.ServiceAccountJSON != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.Firebase.ServiceAccountJSON))
	} else if len(cfg.Firebase.ServiceAccountKey) > 0 {
		opts = append(opts, option.WithCredentialsJSON(cfg.Firebase.ServiceAccountKey))
	}

	firebaseConfig := &firebase.Config{
		ProjectID: cfg.Firebase.ProjectID,
	}

	app, err := firebase.NewApp(ctx, firebaseConfig, opts...)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth client: %v", err)
	}

	// Initialize database client
	var credentialsJSON []byte
	if cfg.Firebase.ServiceAccountJSON != "" {
		credentialsJSON, err = os.ReadFile(cfg.Firebase.ServiceAccountJSON)
		if err != nil {
			log.Fatalf("Failed to read service account file: %v", err)
		}
	} else {
		credentialsJSON = cfg.Firebase.ServiceAccountKey
	}

	dbClient, err := database.NewClient(ctx, cfg.Firebase.ProjectID, credentialsJSON)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer dbClient.Close()

	// Initialize IGDB client
	igdbClient := igdb.NewClient(cfg.IGDB.ClientID, cfg.IGDB.ClientSecret)

	// Start background worker
	backgroundSync := worker.NewBackgroundSync(dbClient, igdbClient, 1*time.Hour)
	go backgroundSync.Start(ctx)
	defer backgroundSync.Stop()

	// Initialize API handler
	apiHandler := api.NewHandler(dbClient)

	// Setup HTTP router
	mux := http.NewServeMux()

	// API routes (with authentication)
	authMiddleware := middleware.AuthMiddleware(authClient)
	mux.Handle("/api/v1/games", authMiddleware(http.HandlerFunc(apiHandler.GetGames)))
	mux.Handle("/api/v1/games/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route to either CreateOrUpdateGame or UpdateGameStatus based on path
		if r.URL.Path == "/api/v1/games" || r.URL.Path == "/api/v1/games/" {
			apiHandler.CreateOrUpdateGame(w, r)
		} else {
			// Check if path ends with /status
			if len(r.URL.Path) > 7 && r.URL.Path[len(r.URL.Path)-7:] == "/status" {
				apiHandler.UpdateGameStatus(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		}
	})))

	// Health check (no authentication)
	mux.HandleFunc("/api/v1/health", apiHandler.HealthCheck)

	// Serve static frontend files
	frontendDistFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("Warning: Failed to load embedded frontend files: %v", err)
		log.Println("API will be available but frontend will not be served")
	} else {
		fileServer := http.FileServer(http.FS(frontendDistFS))
		mux.Handle("/", fileServer)
	}

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down server...")

	// Cancel context to stop background workers
	cancel()

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

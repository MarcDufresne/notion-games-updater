package main

import (
	"context"
	"log"
	"net/http"
	"time"

	gametracker "game-tracker"
	"game-tracker/internal/api"
	"game-tracker/internal/config"
	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
	"game-tracker/internal/middleware"
	"game-tracker/internal/worker"
)

func main() {
	config.Load()
	credentialsFile := config.Get("GOOGLE_APPLICATION_CREDENTIALS")

	ctx := context.Background()
	db, authClient, err := database.NewClient(ctx, credentialsFile)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize IGDB Client
	igdbClientID := config.Get("IGDB_CLIENT_ID")
	igdbClientSecret := config.Get("IGDB_CLIENT_SECRET")
	igdbClient := igdb.NewClient(igdbClientID, igdbClientSecret)

	// Start Background Worker
	bgWorker := worker.NewWorker(db, igdbClient)
	bgWorker.StartBackgroundSync(ctx, 24*time.Hour)

	// Initialize API Handler
	apiHandler := api.NewHandler(db)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	mux := http.NewServeMux()

	// API Router
	apiMux := http.NewServeMux()
	apiHandler.RegisterRoutes(apiMux)

	// Wrap API with Auth Middleware
	mux.Handle("/api/v1/", authMiddleware.Handle(apiMux))

	// Serve Frontend
	mux.Handle("/", http.FileServer(http.FS(gametracker.GetFrontendAssets())))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

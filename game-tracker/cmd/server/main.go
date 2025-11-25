package main

import (
	"context"
	"embed"
	"game-tracker/internal/middleware"
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
	"game-tracker/internal/cache"
	"game-tracker/internal/config"
	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
	"game-tracker/internal/worker"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Firebase Admin SDK
	var opts []option.ClientOption
	if cfg.Firebase.ServiceAccountKey != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(cfg.Firebase.ServiceAccountKey)))
	} else if cfg.Firebase.ServiceAccountJSON != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.Firebase.ServiceAccountJSON))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Get Firebase Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firebase Auth client: %v", err)
	}

	// Initialize Firestore client
	db, err := database.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to Firestore")

	// Initialize IGDB client
	igdbClient := igdb.NewClient(cfg.IGDB.ClientID, cfg.IGDB.ClientSecret)
	log.Println("IGDB client initialized")

	// Initialize search cache (500 entries, 1 hour TTL)
	searchCache := cache.NewCache(500, 1*time.Hour)
	log.Println("Search cache initialized")

	// Start background sync worker
	go worker.StartBackgroundSync(ctx, db, igdbClient)

	// Initialize API handler
	handler := api.NewHandler(db, igdbClient, searchCache, authClient)

	// Setup HTTP server
	mux := http.NewServeMux()

	// Register API routes
	handler.RegisterRoutes(mux)

	// Serve frontend static files
	frontendDist, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("Warning: Failed to load embedded frontend: %v", err)
		log.Println("Frontend static files will not be served")
	} else {
		fileServer := http.FileServer(http.FS(frontendDist))
		mux.Handle("/", fileServer)
		log.Println("Frontend static files configured")
	}

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      middleware.CORS(mux), // Wrap with CORS middleware
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, gracefully shutting down...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		cancel() // Cancel background workers
	}()

	log.Printf("Server starting on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server stopped")
}

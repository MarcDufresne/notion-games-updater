package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jomei/notionapi"
	"github.com/marc/notion-games-updater/internal/config"
	"github.com/marc/notion-games-updater/internal/igdb"
	"github.com/marc/notion-games-updater/internal/usecase"
)

func main() {
	// CLI flags
	runForever := flag.Bool("run-forever", false, "Run the updater in a loop")
	intervalMin := flag.Int("interval", 15, "Interval in minutes between updates (only used with -run-forever)")
	dryRun := flag.Bool("dry-run", false, "Run in dry-run mode (no Notion updates)")
	flag.Parse()

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
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, cancelling operations...")
		cancel()
	}()

	// Initialize clients
	igdbClient := igdb.NewClient(cfg.IGDB.ClientID, cfg.IGDB.ClientSecret)
	notionClient := notionapi.NewClient(notionapi.Token(cfg.Notion.Token))

	// Determine if dry-run is enabled (CLI flag takes precedence)
	isDryRun := *dryRun || cfg.DryRun
	if isDryRun {
		log.Println("Running in DRY-RUN mode - Notion pages will not be updated")
	}

	// Initialize use case
	uc := usecase.NewUpdateGamesUseCase(igdbClient, notionClient, cfg.Notion.DatabaseID, isDryRun)

	// Run the updater
	if *runForever {
		log.Printf("Running updater in loop mode with %d minute interval", *intervalMin)
		for {
			// Check if context is cancelled before running
			if ctx.Err() != nil {
				log.Println("Context cancelled, shutting down...")
				return
			}

			if err := uc.Run(ctx); err != nil {
				if ctx.Err() != nil {
					log.Println("Operation cancelled")
					return
				}
				log.Printf("ERROR: %v", err)
			}

			log.Printf("Sleeping for %d minutes...", *intervalMin)

			// Use context-aware sleep
			select {
			case <-time.After(time.Duration(*intervalMin) * time.Minute):
				// Continue to next iteration
			case <-ctx.Done():
				log.Println("Context cancelled during sleep, shutting down...")
				return
			}
		}
	} else {
		if err := uc.Run(ctx); err != nil {
			if ctx.Err() != nil {
				log.Println("Operation cancelled")
				os.Exit(130) // Standard exit code for SIGINT
			}
			log.Fatalf("ERROR: %v", err)
		}
	}
}

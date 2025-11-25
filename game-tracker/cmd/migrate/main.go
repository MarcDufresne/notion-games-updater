package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"

	"game-tracker/internal/config"
	"game-tracker/internal/database"
	"game-tracker/internal/model"
)

func main() {
	// CLI flags
	userID := flag.String("user-id", "", "Firebase UID for the target user")
	envFile := flag.String("env", "../.env", "Path to parent .env file with Notion credentials")
	flag.Parse()

	// Load parent .env file for Notion credentials
	if err := godotenv.Load(*envFile); err != nil {
		log.Printf("Warning: Could not load %s: %v", *envFile, err)
	}

	// Get user ID from flag or environment variable
	targetUserID := *userID
	if targetUserID == "" {
		targetUserID = os.Getenv("MIGRATION_TARGET_USER_ID")
	}
	if targetUserID == "" {
		log.Fatal("Error: User ID is required. Use --user-id flag or MIGRATION_TARGET_USER_ID env var")
	}

	// Get Notion credentials
	notionToken := os.Getenv("NOTION_TOKEN")
	notionDatabaseID := os.Getenv("NOTION_DATABASE_ID")

	if notionToken == "" || notionDatabaseID == "" {
		log.Fatal("Error: NOTION_TOKEN and NOTION_DATABASE_ID are required in .env file")
	}

	log.Printf("Starting migration from Notion to Firestore for user: %s", targetUserID)

	// Load game-tracker config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Firestore client
	ctx := context.Background()
	db, err := database.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer db.Close()

	log.Println("Connected to Firestore")

	// Initialize Notion client
	notionClient := notionapi.NewClient(notionapi.Token(notionToken))

	// Query all pages from Notion database
	log.Printf("Fetching pages from Notion database: %s", notionDatabaseID)

	pages, err := fetchAllPages(ctx, notionClient, notionDatabaseID)
	if err != nil {
		log.Fatalf("Failed to fetch pages from Notion: %v", err)
	}

	log.Printf("Found %d pages in Notion database", len(pages))

	// Migrate each page
	successCount := 0
	errorCount := 0

	for _, page := range pages {
		game, err := convertNotionPageToGame(page, targetUserID)
		if err != nil {
			log.Printf("ERROR: Failed to convert page %s: %v", page.ID, err)
			errorCount++
			continue
		}

		if err := db.SaveGame(ctx, game); err != nil {
			log.Printf("ERROR: Failed to save game '%s': %v", game.Title, err)
			errorCount++
			continue
		}

		log.Printf("✓ Migrated: %s (Status: %s)", game.Title, game.Status)
		successCount++
	}

	log.Printf("\nMigration complete!")
	log.Printf("Successfully migrated: %d games", successCount)
	log.Printf("Errors: %d", errorCount)
}

func fetchAllPages(ctx context.Context, client *notionapi.Client, databaseID string) ([]*notionapi.Page, error) {
	var allPages []*notionapi.Page
	var cursor notionapi.Cursor

	for {
		resp, err := client.Database.Query(ctx, notionapi.DatabaseID(databaseID), &notionapi.DatabaseQueryRequest{
			StartCursor: cursor,
			PageSize:    100,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query database: %w", err)
		}

		allPages = append(allPages, resp.Results...)

		if !resp.HasMore {
			break
		}
		cursor = resp.NextCursor
	}

	return allPages, nil
}

func convertNotionPageToGame(page *notionapi.Page, userID string) (*model.Game, error) {
	game := &model.Game{
		UserID: userID,
	}

	// Extract title
	if titleProp, ok := page.Properties["Game"]; ok {
		if tp, ok := titleProp.(*notionapi.TitleProperty); ok && len(tp.Title) > 0 {
			game.Title = tp.Title[0].PlainText
		}
	}
	if game.Title == "" {
		return nil, fmt.Errorf("no title found")
	}

	// Extract status
	if statusProp, ok := page.Properties["Status"]; ok {
		if sp, ok := statusProp.(*notionapi.SelectProperty); ok && sp.Select.Name != "" {
			game.Status = mapNotionStatus(sp.Select.Name)
		}
	}
	// Default to Backlog if no status
	if game.Status == "" {
		game.Status = model.StatusBacklog
	}

	// Extract IGDB ID
	if igdbProp, ok := page.Properties["IGDB ID"]; ok {
		if rtp, ok := igdbProp.(*notionapi.RichTextProperty); ok && len(rtp.RichText) > 0 {
			igdbIDStr := strings.TrimPrefix(rtp.RichText[0].PlainText, ":")
			var igdbID int
			if _, err := fmt.Sscanf(igdbIDStr, "%d", &igdbID); err == nil {
				game.IGDBID = &igdbID
			}
		}
	}

	// Extract rating
	if ratingProp, ok := page.Properties["Rating"]; ok {
		if np, ok := ratingProp.(*notionapi.NumberProperty); ok {
			// Notion stores as 0-1, we store as 0-100
			game.Rating = int(np.Number * 100)
		}
	}

	// Extract genres
	if genresProp, ok := page.Properties["Genres"]; ok {
		if msp, ok := genresProp.(*notionapi.MultiSelectProperty); ok {
			for _, option := range msp.MultiSelect {
				game.Genres = append(game.Genres, option.Name)
			}
		}
	}

	// Extract platforms
	if platformsProp, ok := page.Properties["Platforms"]; ok {
		if msp, ok := platformsProp.(*notionapi.MultiSelectProperty); ok {
			for _, option := range msp.MultiSelect {
				game.Platforms = append(game.Platforms, option.Name)
			}
		}
	}

	// Extract release date
	if releaseProp, ok := page.Properties["Release Date"]; ok {
		if dp, ok := releaseProp.(*notionapi.DateProperty); ok && dp.Date != nil {
			releaseDate := dp.Date.Start.Time
			game.ReleaseDate = &releaseDate
		}
	}

	return game, nil
}

func mapNotionStatus(notionStatus string) model.GameStatus {
	// Direct 1:1 mapping
	switch notionStatus {
	case "Backlog":
		return model.StatusBacklog
	case "Break":
		return model.StatusBreak
	case "Playing":
		return model.StatusPlaying
	case "Done":
		return model.StatusDone
	case "Abandoned":
		return model.StatusAbandoned
	case "Won't Play":
		return model.StatusWontPlay
	default:
		log.Printf("Warning: Unknown status '%s', defaulting to Backlog", notionStatus)
		return model.StatusBacklog
	}
}

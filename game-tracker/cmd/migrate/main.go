package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jomei/notionapi"

	"game-tracker/internal/config"
	"game-tracker/internal/database"
	"game-tracker/internal/model"
)

func main() {
	// CLI flags
	notionToken := flag.String("notion-token", "", "Notion API token (or set NOTION_TOKEN env var)")
	notionDatabaseID := flag.String("notion-database-id", "", "Notion database ID (or set NOTION_DATABASE_ID env var)")
	userID := flag.String("user-id", "", "Firebase User ID to assign to migrated games")
	dryRun := flag.Bool("dry-run", false, "Run in dry-run mode (no writes to Firestore)")
	flag.Parse()

	// Get values from flags or environment variables
	token := *notionToken
	if token == "" {
		token = os.Getenv("NOTION_TOKEN")
	}
	dbID := *notionDatabaseID
	if dbID == "" {
		dbID = os.Getenv("NOTION_DATABASE_ID")
	}
	uid := *userID
	if uid == "" {
		uid = os.Getenv("USER_ID")
	}

	// Validate required fields
	if token == "" {
		log.Fatal("Notion token is required (--notion-token or NOTION_TOKEN env var)")
	}
	if dbID == "" {
		log.Fatal("Notion database ID is required (--notion-database-id or NOTION_DATABASE_ID env var)")
	}
	if uid == "" {
		log.Fatal("User ID is required (--user-id or USER_ID env var)")
	}

	// Load Firebase configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Initialize Notion client
	notionClient := notionapi.NewClient(notionapi.Token(token))

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

	var dbClient *database.Client
	if !*dryRun {
		dbClient, err = database.NewClient(ctx, cfg.Firebase.ProjectID, credentialsJSON)
		if err != nil {
			log.Fatalf("Failed to initialize database client: %v", err)
		}
		defer dbClient.Close()
	}

	// Query Notion database
	log.Println("Querying Notion database...")
	query := &notionapi.DatabaseQueryRequest{}
	
	var allPages []notionapi.Page
	hasMore := true
	startCursor := notionapi.Cursor("")

	for hasMore {
		if startCursor != "" {
			query.StartCursor = startCursor
		}

		resp, err := notionClient.Database.Query(ctx, notionapi.DatabaseID(dbID), query)
		if err != nil {
			log.Fatalf("Failed to query Notion database: %v", err)
		}

		allPages = append(allPages, resp.Results...)
		hasMore = resp.HasMore
		if hasMore {
			startCursor = resp.NextCursor
		}
	}

	log.Printf("Found %d pages in Notion database", len(allPages))

	// Migrate each page
	migratedCount := 0
	errorCount := 0

	for i, page := range allPages {
		log.Printf("Processing page %d/%d: %s", i+1, len(allPages), getPageTitle(page))

		game, err := convertNotionPageToGame(page, uid)
		if err != nil {
			log.Printf("ERROR: Failed to convert page: %v", err)
			errorCount++
			continue
		}

		if *dryRun {
			log.Printf("DRY RUN: Would migrate game: %s (Status: %s, IGDB ID: %v)", 
				game.Title, game.Status, game.IGDBID)
		} else {
			// Create a context with user ID for SaveGame
			userCtx := context.WithValue(ctx, "user_id", uid)
			if err := dbClient.SaveGame(userCtx, game); err != nil {
				log.Printf("ERROR: Failed to save game: %v", err)
				errorCount++
				continue
			}
			log.Printf("Migrated: %s", game.Title)
		}

		migratedCount++
	}

	log.Printf("\nMigration complete!")
	log.Printf("Successfully migrated: %d", migratedCount)
	log.Printf("Errors: %d", errorCount)
}

// getPageTitle extracts the title from a Notion page
func getPageTitle(page notionapi.Page) string {
	if titleProp, ok := page.Properties["Game"]; ok {
		if title, ok := titleProp.(*notionapi.TitleProperty); ok && len(title.Title) > 0 {
			return title.Title[0].PlainText
		}
	}
	// Try "Name" as fallback
	if titleProp, ok := page.Properties["Name"]; ok {
		if title, ok := titleProp.(*notionapi.TitleProperty); ok && len(title.Title) > 0 {
			return title.Title[0].PlainText
		}
	}
	return "Unknown"
}

// convertNotionPageToGame converts a Notion page to a Game model
func convertNotionPageToGame(page notionapi.Page, userID string) (*model.Game, error) {
	game := &model.Game{
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Extract title
	game.Title = getPageTitle(page)

	// Extract status
	if statusProp, ok := page.Properties["Status"]; ok {
		if status, ok := statusProp.(*notionapi.SelectProperty); ok && status.Select.Name != "" {
			game.Status = mapNotionStatus(status.Select.Name)
		}
	}
	// Default status if not found
	if game.Status == "" {
		game.Status = model.StatusBacklog
	}

	// Extract IGDB ID
	if igdbProp, ok := page.Properties["IGDB ID"]; ok {
		if igdb, ok := igdbProp.(*notionapi.RichTextProperty); ok && len(igdb.RichText) > 0 {
			igdbID := igdb.RichText[0].PlainText
			// Remove leading colon if present (per the spec)
			igdbID = strings.TrimPrefix(igdbID, ":")
			if igdbID != "" {
				game.IGDBID = &igdbID
			}
		}
	}

	// Extract rating
	if ratingProp, ok := page.Properties["Rating"]; ok {
		if rating, ok := ratingProp.(*notionapi.NumberProperty); ok && rating.Number != 0 {
			game.Rating = int(rating.Number)
		}
	}

	// Extract release date
	if releaseDateProp, ok := page.Properties["Release Date"]; ok {
		if releaseDate, ok := releaseDateProp.(*notionapi.DateProperty); ok && releaseDate.Date != nil {
			if releaseDate.Date.Start != nil {
				startTime := time.Time(*releaseDate.Date.Start)
				game.ReleaseDate = &startTime
			}
		}
	}

	// Extract genres
	if genresProp, ok := page.Properties["Genres"]; ok {
		if genres, ok := genresProp.(*notionapi.MultiSelectProperty); ok {
			for _, genre := range genres.MultiSelect {
				game.Genres = append(game.Genres, genre.Name)
			}
		}
	}

	// Extract platforms
	if platformsProp, ok := page.Properties["Platforms"]; ok {
		if platforms, ok := platformsProp.(*notionapi.MultiSelectProperty); ok {
			for _, platform := range platforms.MultiSelect {
				game.Platforms = append(game.Platforms, platform.Name)
			}
		}
	}

	// Extract cover URL (if exists as URL property)
	if coverProp, ok := page.Properties["Cover URL"]; ok {
		if cover, ok := coverProp.(*notionapi.URLProperty); ok && cover.URL != "" {
			game.CoverURL = cover.URL
		}
	}

	// Set date_played for completed games
	if game.IsComplete() {
		// Try to get from a "Date Played" property if it exists
		if datePlayedProp, ok := page.Properties["Date Played"]; ok {
			if datePlayed, ok := datePlayedProp.(*notionapi.DateProperty); ok && datePlayed.Date != nil {
				if datePlayed.Date.Start != nil {
					startTime := time.Time(*datePlayed.Date.Start)
					game.DatePlayed = &startTime
				}
			}
		}
		// If not found, use the last edited time
		if game.DatePlayed == nil {
			lastEdited := page.LastEditedTime
			game.DatePlayed = &lastEdited
		}
	}

	return game, nil
}

// mapNotionStatus maps Notion status values to Game status values
func mapNotionStatus(notionStatus string) model.GameStatus {
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
		// Default to Backlog if unknown
		return model.StatusBacklog
	}
}

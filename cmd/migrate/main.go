package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"

	"game-tracker/internal/config"
	"game-tracker/internal/database"
	"game-tracker/internal/model"
)

func main() {
	userID := flag.String("user-id", "", "Firebase UID for the target user")
	envFile := flag.String("env", "../.env", "Path to parent .env file with Notion credentials")
	updateMode := flag.Bool("update", false, "Update existing games instead of skipping duplicates")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Printf("Warning: Could not load %s: %v", *envFile, err)
	}

	targetUserID := *userID
	if targetUserID == "" {
		targetUserID = os.Getenv("MIGRATION_TARGET_USER_ID")
	}
	if targetUserID == "" {
		log.Fatal("Error: User ID is required. Use --user-id flag or MIGRATION_TARGET_USER_ID env var")
	}

	notionToken := os.Getenv("NOTION_TOKEN")
	notionDatabaseID := os.Getenv("NOTION_DATABASE_ID")

	if notionToken == "" || notionDatabaseID == "" {
		log.Fatal("Error: NOTION_TOKEN and NOTION_DATABASE_ID are required in .env file")
	}

	log.Printf("Starting migration from Notion to Firestore for user: %s", targetUserID)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	db, err := database.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer db.Close()

	log.Println("Connected to Firestore")

	notionClient := notionapi.NewClient(notionapi.Token(notionToken))

	log.Printf("Fetching pages from Notion database: %s", notionDatabaseID)

	pages, err := fetchAllPages(ctx, notionClient, notionDatabaseID)
	if err != nil {
		log.Fatalf("Failed to fetch pages from Notion: %v", err)
	}

	log.Printf("Found %d pages in Notion database", len(pages))

	successCount := 0
	errorCount := 0
	skippedCount := 0
	updatedCount := 0

	for _, page := range pages {
		game, err := convertNotionPageToGame(page, targetUserID, *debug)
		if err != nil {
			log.Printf("ERROR: Failed to convert page %s: %v", page.ID, err)
			errorCount++
			continue
		}

		if game.IGDBID > 0 {
			existingGame, err := db.GetGameByIGDBID(ctx, targetUserID, game.IGDBID)
			if err != nil {
				log.Printf("ERROR: Failed to check for duplicate '%s': %v", game.Title, err)
				errorCount++
				continue
			}
			if existingGame != nil {
				if *updateMode {
					// Update mode: update the date_played field if it's set in the migration data
					if *debug {
						log.Printf("DEBUG: Checking date_played for '%s': %v (IsZero: %v)", game.Title, game.DatePlayed, game.DatePlayed == nil || game.DatePlayed.IsZero())
					}
					if game.DatePlayed != nil {
						// Check if it's not the epoch sentinel value used by the database
						epoch := time.Unix(0, 0)
						if !game.DatePlayed.Equal(epoch) {
							existingGame.DatePlayed = game.DatePlayed
							if err := db.SaveGame(ctx, existingGame); err != nil {
								log.Printf("ERROR: Failed to update game '%s': %v", game.Title, err)
								errorCount++
								continue
							}
							log.Printf("↻ Updated: %s (Date Played: %s)", existingGame.Title, existingGame.DatePlayed.Format("2006-01-02"))
							updatedCount++
						} else {
							if *debug {
								log.Printf("DEBUG: Skipping '%s' - date is epoch (sentinel value)", game.Title)
							}
							log.Printf("⊘ Skipped: %s (no date played to update)", game.Title)
							skippedCount++
						}
					} else {
						log.Printf("⊘ Skipped: %s (no date played to update)", game.Title)
						skippedCount++
					}
				} else {
					log.Printf("⊘ Skipped duplicate: %s (IGDB ID: %d already exists as '%s')", game.Title, game.IGDBID, existingGame.Title)
					skippedCount++
				}
				continue
			}
		}

		if err := db.SaveGame(ctx, game); err != nil {
			log.Printf("ERROR: Failed to save game '%s': %v", game.Title, err)
			errorCount++
			continue
		}

		datePlayed := "N/A"
		if game.DatePlayed != nil && !game.DatePlayed.IsZero() {
			datePlayed = game.DatePlayed.Format("2006-01-02")
		}
		log.Printf("✓ Migrated: %s (Status: %s, IGDB ID: %d, Date Played: %s) - metadata will be fetched by worker", game.Title, game.Status, game.IGDBID, datePlayed)
		successCount++
	}

	log.Printf("\nMigration complete!")
	log.Printf("Successfully migrated: %d games", successCount)
	if *updateMode {
		log.Printf("Updated existing: %d games", updatedCount)
	}
	log.Printf("Skipped duplicates: %d", skippedCount)
	log.Printf("Errors: %d", errorCount)

	if successCount > 0 {
		log.Printf("\nNote: Game metadata (cover, rating, genres, platforms, release date, URLs)")
		log.Printf("will be automatically fetched by the background worker within 15 minutes.")
	}
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

		// Convert []notionapi.Page to []*notionapi.Page
		for i := range resp.Results {
			allPages = append(allPages, &resp.Results[i])
		}

		if !resp.HasMore {
			break
		}
		cursor = resp.NextCursor
	}

	return allPages, nil
}

func convertNotionPageToGame(page *notionapi.Page, userID string, debug bool) (*model.Game, error) {
	game := &model.Game{
		UserID: userID,
	}

	if debug {
		log.Printf("DEBUG: Converting page %s, available properties: %v", page.ID, getPropertyNames(page))
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
		if debug {
			log.Printf("DEBUG: Found Status property for '%s': %+v", game.Title, statusProp)
		}

		// Try StatusProperty first (newer Notion type)
		if sp, ok := statusProp.(*notionapi.StatusProperty); ok {
			if debug {
				log.Printf("DEBUG: Status StatusProperty for '%s': %+v", game.Title, sp)
				log.Printf("DEBUG: Status Status.Name for '%s': '%s'", game.Title, sp.Status.Name)
			}
			if sp.Status.Name != "" {
				game.Status = mapNotionStatus(sp.Status.Name)
				if debug {
					log.Printf("DEBUG: Mapped status for '%s': '%s' -> '%s'", game.Title, sp.Status.Name, game.Status)
				}
			}
		} else if sp, ok := statusProp.(*notionapi.SelectProperty); ok {
			// Fall back to SelectProperty (older Notion type)
			if debug {
				log.Printf("DEBUG: Status SelectProperty for '%s': %+v", game.Title, sp)
				log.Printf("DEBUG: Status Select.Name for '%s': '%s'", game.Title, sp.Select.Name)
			}
			if sp.Select.Name != "" {
				game.Status = mapNotionStatus(sp.Select.Name)
				if debug {
					log.Printf("DEBUG: Mapped status for '%s': '%s' -> '%s'", game.Title, sp.Select.Name, game.Status)
				}
			}
		} else if debug {
			log.Printf("DEBUG: Status property is neither StatusProperty nor SelectProperty for '%s', type: %T", game.Title, statusProp)
		}
	} else if debug {
		log.Printf("DEBUG: No Status property found for '%s'", game.Title)
	}
	// Default to Backlog if no status
	if game.Status == "" {
		if debug {
			log.Printf("DEBUG: No status set for '%s', defaulting to Backlog", game.Title)
		}
		game.Status = model.StatusBacklog
	}

	// Extract IGDB ID
	if igdbProp, ok := page.Properties["IGDB ID"]; ok {
		if rtp, ok := igdbProp.(*notionapi.RichTextProperty); ok && len(rtp.RichText) > 0 {
			igdbIDStr := strings.TrimPrefix(rtp.RichText[0].PlainText, ":")
			var igdbID int
			if _, err := fmt.Sscanf(igdbIDStr, "%d", &igdbID); err == nil && igdbID > 0 {
				game.IGDBID = igdbID
				game.MatchStatus = model.MatchStatusMatched
			}
		}
	}

	// If no IGDB ID, mark as unmatched
	if game.IGDBID == 0 {
		game.MatchStatus = model.MatchStatusUnmatched
	}

	// NOTE: We don't migrate Rating, Genres, Platforms, Release Date, or URLs
	// These will be automatically fetched by the background worker for games with IGDB IDs

	// Extract date played
	// Try multiple possible property names for date played
	datePropertyNames := []string{"Date Played", "Date played", "date_played", "Played Date", "Played", "Completed Date", "Finished Date"}

	for _, propName := range datePropertyNames {
		if datePlayedProp, ok := page.Properties[propName]; ok {
			if debug {
				log.Printf("DEBUG: Found property '%s' for game '%s'", propName, game.Title)
			}
			if dp, ok := datePlayedProp.(*notionapi.DateProperty); ok {
				if debug {
					log.Printf("DEBUG: DateProperty for '%s': %+v", game.Title, dp)
				}
				if dp.Date != nil && dp.Date.Start != nil {
					// Convert notionapi.Date to time.Time
					t := time.Time(*dp.Date.Start)
					game.DatePlayed = &t
					if debug {
						log.Printf("DEBUG: Successfully parsed date played for '%s': %s (IsZero: %v)", game.Title, t.Format("2006-01-02"), t.IsZero())
					}
					break
				} else if debug {
					log.Printf("DEBUG: DateProperty exists but Date or Date.Start is nil for '%s'", game.Title)
				}
			} else if debug {
				log.Printf("DEBUG: Property '%s' is not a DateProperty for '%s'", propName, game.Title)
			}
		}
	}

	// If no date played found and status is completed, use the page's last_edited_time as fallback
	if game.DatePlayed == nil && (game.Status == model.StatusDone || game.Status == model.StatusAbandoned || game.Status == model.StatusWontPlay) {
		if debug {
			log.Printf("DEBUG: No date played property found for '%s' (completed game), using last_edited_time: %s", game.Title, page.LastEditedTime.Format("2006-01-02"))
		}
		game.DatePlayed = &page.LastEditedTime
	}

	if debug && game.DatePlayed != nil {
		log.Printf("DEBUG: Final DatePlayed for '%s': %s", game.Title, game.DatePlayed.Format("2006-01-02 15:04:05"))
	}

	return game, nil
}

func getPropertyNames(page *notionapi.Page) []string {
	names := make([]string, 0, len(page.Properties))
	for name := range page.Properties {
		names = append(names, name)
	}
	return names
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

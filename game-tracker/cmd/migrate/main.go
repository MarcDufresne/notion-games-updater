package main

import (
	"context"
	"log"
	"strings"
	"time"

	"game-tracker/internal/config"
	"game-tracker/internal/database"
	"game-tracker/internal/model"

	"github.com/jomei/notionapi"
)

func main() {
	config.Load()

	// Initialize Firestore
	credentialsFile := config.Get("GOOGLE_APPLICATION_CREDENTIALS")
	ctx := context.Background()
	db, _, err := database.NewClient(ctx, credentialsFile)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Notion
	notionToken := config.Get("NOTION_TOKEN")
	notionDBID := config.Get("NOTION_DATABASE_ID")
	if notionToken == "" || notionDBID == "" {
		log.Fatal("NOTION_TOKEN and NOTION_DATABASE_ID are required")
	}
	notionClient := notionapi.NewClient(notionapi.Token(notionToken))

	// Iterate and Migrate
	log.Println("Starting migration...")

	hasMore := true
	var startCursor notionapi.Cursor

	for hasMore {
		req := &notionapi.DatabaseQueryRequest{
			StartCursor: startCursor,
		}
		resp, err := notionClient.Database.Query(ctx, notionapi.DatabaseID(notionDBID), req)
		if err != nil {
			log.Fatalf("Failed to query Notion: %v", err)
		}

		for _, page := range resp.Results {
			game := mapNotionPageToGame(page)
			if game != nil {
				// Hardcode user ID for migration, or pass as arg.
				// Assuming single user for now as per "Personal game library manager".
				// Or we can use an env var MIGRATION_USER_ID.
				game.UserID = config.Get("MIGRATION_USER_ID")
				if game.UserID == "" {
					log.Fatal("MIGRATION_USER_ID env var is required")
				}

				if err := db.SaveGame(ctx, game); err != nil {
					log.Printf("Failed to save game %s: %v", game.Title, err)
				} else {
					log.Printf("Migrated: %s", game.Title)
				}
			}
		}

		hasMore = resp.HasMore
		startCursor = resp.NextCursor
	}

	log.Println("Migration complete.")
}

func mapNotionPageToGame(page notionapi.Page) *model.Game {
	props := page.Properties

	titleList, ok := props["Name"].(*notionapi.TitleProperty)
	if !ok || len(titleList.Title) == 0 {
		return nil
	}
	title := titleList.Title[0].PlainText

	status := model.StatusBacklog // Default
	if statusProp, ok := props["Status"].(*notionapi.SelectProperty); ok {
		switch statusProp.Select.Name {
		case "Backlog": status = model.StatusBacklog
		case "Break": status = model.StatusBreak
		case "Playing": status = model.StatusPlaying
		case "Done": status = model.StatusDone
		case "Abandoned": status = model.StatusAbandoned
		case "Won't Play": status = model.StatusWontPlay
		}
	}

	var igdbID *string
	if idProp, ok := props["IGDB ID"].(*notionapi.RichTextProperty); ok && len(idProp.RichText) > 0 {
		val := strings.TrimPrefix(idProp.RichText[0].PlainText, ":")
		igdbID = &val
	}

	// Parse other fields as needed based on "Data Model (Firestore)" description
	// and "Mapping" in instructions.

	return &model.Game{
		Title:     title,
		Status:    status,
		IGDBID:    igdbID,
		CreatedAt: time.Now(), // Notion doesn't easily give created time without extra lookups, or use page.CreatedTime
		UpdatedAt: time.Now(),
	}
}

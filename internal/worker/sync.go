package worker

import (
	"context"
	"log"
	"time"

	"game-tracker/internal/api"
	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
	"game-tracker/internal/model"
)

// StartBackgroundSync starts a background goroutine that syncs game metadata from IGDB every 15 minutes
func StartBackgroundSync(ctx context.Context, db *database.Client, igdbClient *igdb.Client) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	log.Println("Starting background sync worker (15 minute interval)")

	// Run initial sync
	syncGames(ctx, db, igdbClient)

	for {
		select {
		case <-ctx.Done():
			log.Println("Background sync worker stopped")
			return
		case <-ticker.C:
			syncGames(ctx, db, igdbClient)
		}
	}
}

func syncGames(ctx context.Context, db *database.Client, igdbClient *igdb.Client) {
	log.Println("Starting background game metadata sync...")

	syncMatchedGames(ctx, db, igdbClient)
	matchUnmatchedGames(ctx, db, igdbClient)
}

func syncMatchedGames(ctx context.Context, db *database.Client, igdbClient *igdb.Client) {
	games, err := db.GetGamesWithIGDBID(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to fetch games for sync: %v", err)
		return
	}

	log.Printf("Found %d matched games to sync", len(games))

	successCount := 0
	errorCount := 0

	for _, game := range games {
		if game.IGDBID == 0 {
			continue
		}

		igdbGame, err := igdbClient.GetGameByID(game.IGDBID)
		if err != nil {
			log.Printf("ERROR: Failed to fetch IGDB data for game '%s' (ID: %d): %v", game.Title, game.IGDBID, err)

			game.LastSyncError = err.Error()
			if saveErr := db.SaveGame(ctx, game); saveErr != nil {
				log.Printf("ERROR: Failed to update sync error for game '%s': %v", game.Title, saveErr)
			}
			errorCount++
			continue
		}

		updated := api.EnrichGameFromIGDB(game, igdbGame, true)

		if updated {
			game.LastSyncError = ""

			if err := db.SaveGame(ctx, game); err != nil {
				log.Printf("ERROR: Failed to save updated game '%s': %v", game.Title, err)
				errorCount++
				continue
			}

			log.Printf("Successfully synced game: %s", game.Title)
			successCount++
		}
	}

	log.Printf("Background sync complete: %d successful, %d errors", successCount, errorCount)
}

func matchUnmatchedGames(ctx context.Context, db *database.Client, igdbClient *igdb.Client) {
	games, err := db.GetUnmatchedGames(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to fetch unmatched games: %v", err)
		return
	}

	if len(games) == 0 {
		return
	}

	log.Printf("Found %d unmatched games to process", len(games))

	matchedCount := 0
	multipleCount := 0
	noMatchCount := 0

	for _, game := range games {
		// Skip if already marked as needs review
		if game.MatchStatus == model.MatchStatusNeedsReview {
			continue
		}

		// Search IGDB for this game title
		searchResults, err := igdbClient.Search(game.Title)
		if err != nil {
			log.Printf("ERROR: Failed to search IGDB for '%s': %v", game.Title, err)
			continue
		}

		// Handle different match scenarios
		if len(searchResults) == 0 {
			// No matches found
			log.Printf("No IGDB matches found for: %s", game.Title)
			game.MatchStatus = model.MatchStatusNoMatch
			if saveErr := db.SaveGame(ctx, game); saveErr != nil {
				log.Printf("ERROR: Failed to update match status for '%s': %v", game.Title, saveErr)
			}
			noMatchCount++
		} else if len(searchResults) == 1 {
			// Single match - check for duplicates before automatically linking
			igdbID := searchResults[0].ID

			// Check if this IGDB ID is already used by another game for this user
			existingGame, err := db.GetGameByIGDBID(ctx, game.UserID, igdbID)
			if err != nil {
				log.Printf("ERROR: Failed to check for duplicate IGDB ID %d for '%s': %v", igdbID, game.Title, err)
				continue
			}

			if existingGame != nil && existingGame.ID != game.ID {
				// Another game already has this IGDB ID - mark as needs review
				log.Printf("IGDB ID %d already used by '%s' - marking '%s' as needs review", igdbID, existingGame.Title, game.Title)
				game.MatchStatus = model.MatchStatusNeedsReview
				if saveErr := db.SaveGame(ctx, game); saveErr != nil {
					log.Printf("ERROR: Failed to update match status for '%s': %v", game.Title, saveErr)
				}
				multipleCount++ // Count as needing review
				continue
			}

			// No duplicate found - fetch details and auto-match
			igdbGame, err := igdbClient.GetGameByID(igdbID)
			if err != nil {
				log.Printf("ERROR: Failed to fetch IGDB details for '%s' (ID: %d): %v", game.Title, igdbID, err)
				continue
			}

			game.IGDBID = igdbID
			game.MatchStatus = model.MatchStatusMatched
			api.EnrichGameFromIGDB(game, igdbGame, false)

			if saveErr := db.SaveGame(ctx, game); saveErr != nil {
				log.Printf("ERROR: Failed to save matched game '%s': %v", game.Title, saveErr)
				continue
			}

			log.Printf("Automatically matched: %s -> IGDB ID: %d", game.Title, igdbID)
			matchedCount++
		} else {
			// Multiple matches - mark for user review
			log.Printf("Multiple IGDB matches found for '%s' (%d results) - marking for review", game.Title, len(searchResults))
			game.MatchStatus = model.MatchStatusMultiple
			if saveErr := db.SaveGame(ctx, game); saveErr != nil {
				log.Printf("ERROR: Failed to update match status for '%s': %v", game.Title, saveErr)
			}
			multipleCount++
		}
	}

	if matchedCount > 0 || multipleCount > 0 || noMatchCount > 0 {
		log.Printf("Unmatched game processing complete: %d auto-matched, %d multiple matches, %d no match", matchedCount, multipleCount, noMatchCount)
	}
}

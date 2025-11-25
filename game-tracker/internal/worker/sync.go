package worker

import (
	"context"
	"game-tracker/internal/legacy_domain"
	"log"
	"time"

	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
	"game-tracker/internal/model"
)

// StartBackgroundSync starts a background goroutine that syncs game metadata from IGDB every 15 minutes
func StartBackgroundSync(ctx context.Context, db *database.Client, igdbClient *igdb.Client) {
	ticker := time.NewTicker(15 * time.Minute)
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

	// Get all games with IGDB IDs
	games, err := db.GetGamesWithIGDBID(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to fetch games for sync: %v", err)
		return
	}

	log.Printf("Found %d games to sync", len(games))

	successCount := 0
	errorCount := 0

	for _, game := range games {
		if game.IGDBID == nil {
			continue
		}

		// Fetch latest metadata from IGDB
		igdbGame, err := igdbClient.GetGameByID(*game.IGDBID)
		if err != nil {
			log.Printf("ERROR: Failed to fetch IGDB data for game '%s' (ID: %d): %v", game.Title, *game.IGDBID, err)

			// Update last_sync_error field
			game.LastSyncError = err.Error()
			if saveErr := db.SaveGame(ctx, game); saveErr != nil {
				log.Printf("ERROR: Failed to update sync error for game '%s': %v", game.Title, saveErr)
			}
			errorCount++
			continue
		}

		// Update game metadata
		updated := updateGameFromIGDB(game, igdbGame)

		if updated {
			// Clear any previous sync errors
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

// updateGameFromIGDB updates game fields with fresh IGDB data, returns true if any changes were made
func updateGameFromIGDB(game *model.Game, igdbGame *legacy_domain.Game) bool {
	updated := false

	// Update cover URL
	if igdbGame.Cover != nil {
		newCoverURL := igdbGame.Cover.CoverBig2xURL()
		if game.CoverURL != newCoverURL {
			game.CoverURL = newCoverURL
			updated = true
		}
	}

	// Update rating
	if igdbGame.AggregatedRating != nil {
		newRating := int(*igdbGame.AggregatedRating)
		if game.Rating != newRating {
			game.Rating = newRating
			updated = true
		}
	}

	// Update genres
	if len(igdbGame.Genres) > 0 {
		newGenres := make([]string, len(igdbGame.Genres))
		for i, genre := range igdbGame.Genres {
			newGenres[i] = genre.Name
		}
		if !stringSlicesEqual(game.Genres, newGenres) {
			game.Genres = newGenres
			updated = true
		}
	}

	// Update platforms
	if len(igdbGame.Platforms) > 0 {
		newPlatforms := make([]string, 0, len(igdbGame.Platforms))
		for _, platform := range igdbGame.Platforms {
			if platform.Abbreviation == "" {
				continue
			}
			// Skip Stadia
			if platform.Abbreviation == "Stadia" {
				continue
			}
			newPlatforms = append(newPlatforms, platform.Abbreviation)
		}
		if !stringSlicesEqual(game.Platforms, newPlatforms) {
			game.Platforms = newPlatforms
			updated = true
		}
	}

	// Update release date
	if igdbGame.FirstReleaseDate != nil {
		releaseDate := time.Unix(*igdbGame.FirstReleaseDate, 0)
		if game.ReleaseDate == nil || !game.ReleaseDate.Equal(releaseDate) {
			game.ReleaseDate = &releaseDate
			updated = true
		}
	}

	return updated
}

// stringSlicesEqual checks if two string slices are equal
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

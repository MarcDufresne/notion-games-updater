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

	// Sync matched games (with IGDB IDs)
	syncMatchedGames(ctx, db, igdbClient)

	// Try to match unmatched games
	matchUnmatchedGames(ctx, db, igdbClient)
}

func syncMatchedGames(ctx context.Context, db *database.Client, igdbClient *igdb.Client) {
	// Get all games with IGDB IDs
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

		// Fetch latest metadata from IGDB
		igdbGame, err := igdbClient.GetGameByID(game.IGDBID)
		if err != nil {
			log.Printf("ERROR: Failed to fetch IGDB data for game '%s' (ID: %d): %v", game.Title, game.IGDBID, err)

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
			updateGameFromIGDB(game, igdbGame)

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

// updateGameFromIGDB updates game fields with fresh IGDB data, returns true if any changes were made
func updateGameFromIGDB(game *model.Game, igdbGame *legacy_domain.Game) bool {
	updated := false

	// Update title
	if game.Title != igdbGame.Name {
		game.Title = igdbGame.Name
		updated = true
	}

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

	// Update Steam URL and Official URL
	if len(igdbGame.Websites) > 0 {
		newSteamURL := ""
		newOfficialURL := ""
		for _, website := range igdbGame.Websites {
			if website.Type == legacy_domain.WebsiteCategorySteam {
				newSteamURL = website.URL
			} else if website.Type == legacy_domain.WebsiteCategoryOfficial {
				newOfficialURL = website.URL
			}
		}
		if game.SteamURL != newSteamURL {
			game.SteamURL = newSteamURL
			updated = true
		}
		if game.OfficialURL != newOfficialURL {
			game.OfficialURL = newOfficialURL
			updated = true
		}
	}

	// Clear any previous sync errors since we successfully fetched data
	if game.LastSyncError != "" {
		game.LastSyncError = ""
		updated = true
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

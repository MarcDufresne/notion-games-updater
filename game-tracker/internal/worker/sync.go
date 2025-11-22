package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
	"game-tracker/internal/legacy_domain"
)

type BackgroundSync struct {
	db         *database.Client
	igdbClient *igdb.Client
	interval   time.Duration
	stopChan   chan struct{}
}

// NewBackgroundSync creates a new background sync worker
func NewBackgroundSync(db *database.Client, igdbClient *igdb.Client, interval time.Duration) *BackgroundSync {
	return &BackgroundSync{
		db:         db,
		igdbClient: igdbClient,
		interval:   interval,
		stopChan:   make(chan struct{}),
	}
}

// Start begins the background sync loop
func (bs *BackgroundSync) Start(ctx context.Context) {
	ticker := time.NewTicker(bs.interval)
	defer ticker.Stop()

	// Run once immediately
	bs.runSync(ctx)

	for {
		select {
		case <-ticker.C:
			bs.runSync(ctx)
		case <-bs.stopChan:
			log.Println("Background sync stopped")
			return
		case <-ctx.Done():
			log.Println("Background sync stopped due to context cancellation")
			return
		}
	}
}

// Stop stops the background sync
func (bs *BackgroundSync) Stop() {
	close(bs.stopChan)
}

// runSync executes a sync cycle
func (bs *BackgroundSync) runSync(ctx context.Context) {
	log.Println("Starting background sync...")

	// Update games with IGDB IDs
	if err := bs.updateGamesWithIGDBID(ctx); err != nil {
		log.Printf("Error updating games with IGDB ID: %v", err)
	}

	// Search for games without IGDB IDs
	if err := bs.searchGamesWithoutIGDBID(ctx); err != nil {
		log.Printf("Error searching games without IGDB ID: %v", err)
	}

	log.Println("Background sync completed")
}

// updateGamesWithIGDBID updates metadata for games that have an IGDB ID
func (bs *BackgroundSync) updateGamesWithIGDBID(ctx context.Context) error {
	games, err := bs.db.GetGamesWithIGDBID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get games with IGDB ID: %w", err)
	}

	log.Printf("Found %d games with IGDB ID to update", len(games))

	for _, game := range games {
		if game.IGDBID == nil || *game.IGDBID == "" {
			continue
		}

		// Fetch game data from IGDB
		query := fmt.Sprintf(`
			fields name, cover.image_id, aggregated_rating, first_release_date, 
			       genres.name, platforms.name, platforms.abbreviation;
			where id = %s;
		`, *game.IGDBID)

		body, err := bs.igdbClient.Request("games", query)
		if err != nil {
			log.Printf("Error fetching IGDB data for game %s (IGDB ID: %s): %v", game.Title, *game.IGDBID, err)
			continue
		}

		var igdbGames []legacy_domain.Game
		if err := json.Unmarshal(body, &igdbGames); err != nil {
			log.Printf("Error parsing IGDB data for game %s: %v", game.Title, err)
			continue
		}

		if len(igdbGames) == 0 {
			log.Printf("No IGDB data found for game %s (IGDB ID: %s)", game.Title, *game.IGDBID)
			continue
		}

		igdbGame := igdbGames[0]

		// Update game metadata
		updated := false

		if igdbGame.Cover != nil {
			newCoverURL := igdbGame.Cover.CoverBig2xURL()
			if game.CoverURL != newCoverURL {
				game.CoverURL = newCoverURL
				updated = true
			}
		}

		if igdbGame.AggregatedRating != nil {
			newRating := int(*igdbGame.AggregatedRating)
			if game.Rating != newRating {
				game.Rating = newRating
				updated = true
			}
		}

		if igdbGame.FirstReleaseDate != nil {
			releaseTime := time.Unix(*igdbGame.FirstReleaseDate, 0)
			if game.ReleaseDate == nil || !game.ReleaseDate.Equal(releaseTime) {
				game.ReleaseDate = &releaseTime
				updated = true
			}
		}

		if len(igdbGame.Genres) > 0 {
			var genres []string
			for _, genre := range igdbGame.Genres {
				genres = append(genres, genre.Name)
			}
			// Simple comparison - update if different
			if fmt.Sprintf("%v", game.Genres) != fmt.Sprintf("%v", genres) {
				game.Genres = genres
				updated = true
			}
		}

		if len(igdbGame.Platforms) > 0 {
			var platforms []string
			for _, platform := range igdbGame.Platforms {
				platforms = append(platforms, platform.Name)
			}
			// Simple comparison - update if different
			if fmt.Sprintf("%v", game.Platforms) != fmt.Sprintf("%v", platforms) {
				game.Platforms = platforms
				updated = true
			}
		}

		if updated {
			// Create a context with the user ID for SaveGame
			userCtx := context.WithValue(ctx, "user_id", game.UserID)
			if err := bs.db.SaveGame(userCtx, game); err != nil {
				log.Printf("Error saving updated game %s: %v", game.Title, err)
				continue
			}
			log.Printf("Updated game: %s", game.Title)
		}
	}

	return nil
}

// searchGamesWithoutIGDBID searches IGDB for games without an IGDB ID
func (bs *BackgroundSync) searchGamesWithoutIGDBID(ctx context.Context) error {
	games, err := bs.db.GetGamesWithoutIGDBID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get games without IGDB ID: %w", err)
	}

	log.Printf("Found %d games without IGDB ID to search", len(games))

	for _, game := range games {
		// Search IGDB by title
		query := fmt.Sprintf(`
			search "%s";
			fields game.id, game.name, game.cover.image_id;
			limit 5;
		`, game.Title)

		body, err := bs.igdbClient.Request("search", query)
		if err != nil {
			log.Printf("Error searching IGDB for game %s: %v", game.Title, err)
			continue
		}

		var searchResults []legacy_domain.SearchResult
		if err := json.Unmarshal(body, &searchResults); err != nil {
			log.Printf("Error parsing IGDB search results for game %s: %v", game.Title, err)
			continue
		}

		if len(searchResults) == 0 {
			log.Printf("No IGDB results found for game: %s", game.Title)
			continue
		}

		// Log potential matches
		log.Printf("Found %d potential matches for game '%s':", len(searchResults), game.Title)
		for i, result := range searchResults {
			if result.Game != nil {
				log.Printf("  %d. %s (IGDB ID: %d)", i+1, result.Game.Name, result.Game.ID)
			}
		}
		log.Printf("  Manual linking required for game: %s", game.Title)
	}

	return nil
}

// RunOnce runs a single sync cycle (useful for testing)
func (bs *BackgroundSync) RunOnce(ctx context.Context) error {
	bs.runSync(ctx)
	return nil
}

// Helper function to convert int to string pointer
func intToStringPtr(i int) *string {
	s := strconv.Itoa(i)
	return &s
}

package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
)

type Worker struct {
	db   *database.Client
	igdb *igdb.Client
}

func NewWorker(db *database.Client, igdbClient *igdb.Client) *Worker {
	return &Worker{
		db:   db,
		igdb: igdbClient,
	}
}

func (w *Worker) StartBackgroundSync(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := w.Sync(ctx); err != nil {
					log.Printf("Error during background sync: %v", err)
				}
			}
		}
	}()
}

func (w *Worker) Sync(ctx context.Context) error {
	log.Println("Running background sync...")

	games, err := w.db.GetAllGames(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all games: %w", err)
	}

	for _, game := range games {
		if game.IGDBID != nil && *game.IGDBID != "" {
			// Update game logic here
			// For now, just log
			// log.Printf("Syncing game: %s (IGDB ID: %s)", game.Title, *game.IGDBID)
		} else {
			// Search IGDB logic here
			// log.Printf("Searching IGDB for game: %s", game.Title)
			// query := fmt.Sprintf("search \"%s\"; fields name, cover.image_id, first_release_date, platforms.name; limit 5;", game.Title)
			// _ = query
		}
	}

	return nil
}

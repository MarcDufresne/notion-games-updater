package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"game-tracker/internal/config"
	"game-tracker/internal/model"
)

const gamesCollection = "games"

type Client struct {
	firestore *firestore.Client
}

// NewClient creates a new Firestore client with credential support for both raw JSON and file path
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	var opts []option.ClientOption

	// Precedence: raw JSON key > file path
	if cfg.Firebase.ServiceAccountKey != "" {
		// Use raw JSON key
		opts = append(opts, option.WithCredentialsJSON([]byte(cfg.Firebase.ServiceAccountKey)))
	} else if cfg.Firebase.ServiceAccountJSON != "" {
		// Check if file exists
		if _, err := os.Stat(cfg.Firebase.ServiceAccountJSON); err != nil {
			return nil, fmt.Errorf("service account file not found: %w", err)
		}
		opts = append(opts, option.WithCredentialsFile(cfg.Firebase.ServiceAccountJSON))
	} else {
		return nil, fmt.Errorf("no Firebase credentials provided")
	}

	// Initialize Firebase app
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Get Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return &Client{
		firestore: firestoreClient,
	}, nil
}

// Close closes the Firestore client
func (c *Client) Close() error {
	return c.firestore.Close()
}

// SaveGame creates or updates a game, enforcing user_id from context
func (c *Client) SaveGame(ctx context.Context, game *model.Game) error {
	now := time.Now()

	// Set timestamps
	if game.CreatedAt.IsZero() {
		game.CreatedAt = now
	}
	game.UpdatedAt = now

	// Set sentinel values for nil dates to enable proper sorting in Firestore
	// Games without release dates get far future date (sort to end)
	if game.ReleaseDate == nil {
		farFuture := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
		game.ReleaseDate = &farFuture
	}

	// Games without date_played get epoch (sort to end when descending)
	if game.DatePlayed == nil {
		epoch := time.Unix(0, 0)
		game.DatePlayed = &epoch
	}

	// Set match status if not already set
	if game.MatchStatus == "" {
		if game.IGDBID > 0 {
			game.MatchStatus = model.MatchStatusMatched
		} else {
			game.MatchStatus = model.MatchStatusUnmatched
		}
	}

	// Generate ID if not present
	if game.ID == "" {
		docRef := c.firestore.Collection(gamesCollection).NewDoc()
		game.ID = docRef.ID
	}

	// Save to Firestore
	_, err := c.firestore.Collection(gamesCollection).Doc(game.ID).Set(ctx, game)
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}

	return nil
}

// GetGame retrieves a single game by ID
func (c *Client) GetGame(ctx context.Context, gameID string) (*model.Game, error) {
	doc, err := c.firestore.Collection(gamesCollection).Doc(gameID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	var game model.Game
	if err := doc.DataTo(&game); err != nil {
		return nil, fmt.Errorf("failed to parse game: %w", err)
	}

	return &game, nil
}

func (c *Client) GetGameByIGDBID(ctx context.Context, userID string, igdbID int) (*model.Game, error) {
	docs, err := c.firestore.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("igdb_id", "==", igdbID).
		Limit(1).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query game by IGDB ID: %w", err)
	}

	if len(docs) == 0 {
		return nil, nil
	}

	var game model.Game
	if err := docs[0].DataTo(&game); err != nil {
		return nil, fmt.Errorf("failed to parse game: %w", err)
	}

	return &game, nil
}

// GetGames retrieves games for a user with optional status filter
func (c *Client) GetGames(ctx context.Context, userID string, statuses ...model.GameStatus) ([]*model.Game, error) {
	query := c.firestore.Collection(gamesCollection).Where("user_id", "==", userID)

	if len(statuses) > 0 {
		query = query.Where("status", "in", statusesToInterfaces(statuses))
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query games: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetBacklog retrieves backlog games (Backlog and Break) sorted by release date ASC
// Games without release dates have sentinel far-future date and sort to end
func (c *Client) GetBacklog(ctx context.Context, userID string) ([]*model.Game, error) {
	docs, err := c.firestore.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("status", "in", []interface{}{model.StatusBacklog, model.StatusBreak}).
		OrderBy("release_date", firestore.Asc).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query backlog: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetUpcoming retrieves backlog games with release dates from the last month onwards, sorted by release date ASC
func (c *Client) GetUpcoming(ctx context.Context, userID string) ([]*model.Game, error) {
	// Calculate date one month ago
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	docs, err := c.firestore.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("status", "in", []interface{}{model.StatusBacklog, model.StatusBreak}).
		Where("release_date", ">=", oneMonthAgo).
		OrderBy("release_date", firestore.Asc).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query upcoming games: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetPlaying retrieves currently playing games sorted by updated_at DESC
func (c *Client) GetPlaying(ctx context.Context, userID string) ([]*model.Game, error) {
	docs, err := c.firestore.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("status", "==", model.StatusPlaying).
		OrderBy("updated_at", firestore.Desc).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query playing: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetHistory retrieves completed games (Done, Abandoned, Won't Play) sorted by date_played DESC
// Games without date_played have sentinel epoch date and sort to end
func (c *Client) GetHistory(ctx context.Context, userID string) ([]*model.Game, error) {
	docs, err := c.firestore.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("status", "in", []interface{}{model.StatusDone, model.StatusAbandoned, model.StatusWontPlay}).
		OrderBy("date_played", firestore.Desc).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetAllGames retrieves all games for a user sorted by release date DESC (newest first)
// Games without release dates (sentinel values) sort to the end
func (c *Client) GetAllGames(ctx context.Context, userID string) ([]*model.Game, error) {
	docs, err := c.firestore.Collection(gamesCollection).
		Where("user_id", "==", userID).
		OrderBy("release_date", firestore.Desc).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query all games: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetGamesWithIGDBID retrieves all games that have an IGDB ID for background sync
func (c *Client) GetGamesWithIGDBID(ctx context.Context) ([]*model.Game, error) {
	docs, err := c.firestore.Collection(gamesCollection).
		Where("igdb_id", ">", 0).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query games with IGDB ID: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetUnmatchedGames retrieves games without IGDB IDs (manual entries needing matching)
func (c *Client) GetUnmatchedGames(ctx context.Context) ([]*model.Game, error) {
	docs, err := c.firestore.Collection(gamesCollection).
		Where("igdb_id", "==", 0).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to query unmatched games: %w", err)
	}

	games := make([]*model.Game, 0, len(docs))
	for _, doc := range docs {
		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

// UpdateGameStatus updates only the status of a game
func (c *Client) UpdateGameStatus(ctx context.Context, gameID string, status model.GameStatus, datePlayed *time.Time) error {
	updates := []firestore.Update{
		{Path: "status", Value: status},
		{Path: "updated_at", Value: time.Now()},
	}

	// If status is being changed to a completed state, update date_played
	if status == model.StatusDone || status == model.StatusAbandoned || status == model.StatusWontPlay {
		// Use provided date or default to now
		playedDate := time.Now()
		if datePlayed != nil {
			playedDate = *datePlayed
			log.Printf("DEBUG: Using provided date_played: %v", playedDate)
		} else {
			log.Printf("DEBUG: No date_played provided, using current time: %v", playedDate)
		}
		updates = append(updates, firestore.Update{Path: "date_played", Value: playedDate})
	}

	_, err := c.firestore.Collection(gamesCollection).Doc(gameID).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to update game status: %w", err)
	}

	return nil
}

// DeleteGame permanently deletes a game from the database
func (c *Client) DeleteGame(ctx context.Context, gameID string) error {
	_, err := c.firestore.Collection(gamesCollection).Doc(gameID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	log.Printf("Successfully deleted game from Firestore: %s", gameID)
	return nil
}

// Helper function to convert statuses to interface slice for Firestore
func statusesToInterfaces(statuses []model.GameStatus) []interface{} {
	result := make([]interface{}, len(statuses))
	for i, s := range statuses {
		result[i] = s
	}
	return result
}

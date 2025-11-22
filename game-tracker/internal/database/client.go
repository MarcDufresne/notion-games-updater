package database

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"game-tracker/internal/model"
)

const (
	gamesCollection = "games"
)

type Client struct {
	firestoreClient *firestore.Client
}

// NewClient creates a new Firestore database client
func NewClient(ctx context.Context, projectID string, credentialsJSON []byte) (*Client, error) {
	var opts []option.ClientOption
	if len(credentialsJSON) > 0 {
		opts = append(opts, option.WithCredentialsJSON(credentialsJSON))
	}

	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore client: %w", err)
	}

	return &Client{
		firestoreClient: client,
	}, nil
}

// Close closes the Firestore client
func (c *Client) Close() error {
	return c.firestoreClient.Close()
}

// SaveGame saves or updates a game in Firestore
// The user_id from context will overwrite the input user_id for security
func (c *Client) SaveGame(ctx context.Context, game *model.Game) error {
	// Extract user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return fmt.Errorf("user_id not found in context")
	}

	// Enforce user_id from context
	game.UserID = userID

	// Set timestamps
	now := time.Now()
	if game.ID == "" {
		// Generate new ID for new games
		game.ID = c.firestoreClient.Collection(gamesCollection).NewDoc().ID
		game.CreatedAt = now
	}
	game.UpdatedAt = now

	// Save to Firestore
	_, err := c.firestoreClient.Collection(gamesCollection).Doc(game.ID).Set(ctx, game)
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}

	return nil
}

// GetGame retrieves a single game by ID
func (c *Client) GetGame(ctx context.Context, userID, gameID string) (*model.Game, error) {
	doc, err := c.firestoreClient.Collection(gamesCollection).Doc(gameID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	var game model.Game
	if err := doc.DataTo(&game); err != nil {
		return nil, fmt.Errorf("failed to parse game data: %w", err)
	}

	// Verify the game belongs to the user
	if game.UserID != userID {
		return nil, fmt.Errorf("game not found")
	}

	return &game, nil
}

// GetGames retrieves games for a user with optional filtering
func (c *Client) GetGames(ctx context.Context, userID string, statuses []model.GameStatus) ([]*model.Game, error) {
	query := c.firestoreClient.Collection(gamesCollection).Where("user_id", "==", userID)

	// Apply status filter if provided
	if len(statuses) > 0 {
		query = query.Where("status", "in", statuses)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var games []*model.Game
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate games: %w", err)
		}

		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game data: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}

// GetBacklog retrieves games in the backlog view (Backlog or Break status)
// Sorted by release_date ASC
func (c *Client) GetBacklog(ctx context.Context, userID string) ([]*model.Game, error) {
	query := c.firestoreClient.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("status", "in", []interface{}{model.StatusBacklog, model.StatusBreak}).
		OrderBy("release_date", firestore.Asc)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var games []*model.Game
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate backlog games: %w", err)
		}

		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game data: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}

// GetPlaying retrieves games currently being played
// Sorted by updated_at DESC
func (c *Client) GetPlaying(ctx context.Context, userID string) ([]*model.Game, error) {
	query := c.firestoreClient.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("status", "==", model.StatusPlaying).
		OrderBy("updated_at", firestore.Desc)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var games []*model.Game
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate playing games: %w", err)
		}

		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game data: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}

// GetHistory retrieves completed games (Done, Abandoned, Won't Play)
// Sorted by date_played DESC
func (c *Client) GetHistory(ctx context.Context, userID string) ([]*model.Game, error) {
	query := c.firestoreClient.Collection(gamesCollection).
		Where("user_id", "==", userID).
		Where("status", "in", []interface{}{model.StatusDone, model.StatusAbandoned, model.StatusWontPlay}).
		OrderBy("date_played", firestore.Desc)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var games []*model.Game
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate history games: %w", err)
		}

		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game data: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}

// GetGamesWithIGDBID retrieves all games that have an IGDB ID for sync
func (c *Client) GetGamesWithIGDBID(ctx context.Context) ([]*model.Game, error) {
	query := c.firestoreClient.Collection(gamesCollection).
		Where("igdb_id", "!=", nil)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var games []*model.Game
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate games with IGDB ID: %w", err)
		}

		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game data: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}

// GetGamesWithoutIGDBID retrieves all games without an IGDB ID for search
func (c *Client) GetGamesWithoutIGDBID(ctx context.Context) ([]*model.Game, error) {
	query := c.firestoreClient.Collection(gamesCollection).
		Where("igdb_id", "==", nil)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var games []*model.Game
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate games without IGDB ID: %w", err)
		}

		var game model.Game
		if err := doc.DataTo(&game); err != nil {
			return nil, fmt.Errorf("failed to parse game data: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}

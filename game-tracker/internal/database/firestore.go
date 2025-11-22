package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"game-tracker/internal/model"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	client *firestore.Client
	auth   *auth.Client
}

func NewClient(ctx context.Context, credentialsFile string) (*Client, *auth.Client, error) {
	var opts []option.ClientOption
	if credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing firestore: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing auth: %v", err)
	}

	return &Client{client: client, auth: authClient}, authClient, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) SaveGame(ctx context.Context, game *model.Game) error {
	if game.ID == "" {
		// Auto-generate ID if not provided
		ref := c.client.Collection("games").NewDoc()
		game.ID = ref.ID
		game.CreatedAt = time.Now()
	}
	game.UpdatedAt = time.Now()

	_, err := c.client.Collection("games").Doc(game.ID).Set(ctx, game)
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}
	return nil
}

func (c *Client) GetGames(ctx context.Context, userID string, view string) ([]model.Game, error) {
	col := c.client.Collection("games")
	query := col.Where("user_id", "==", userID)

	switch view {
	case "backlog":
		query = query.Where("status", "in", []model.GameStatus{model.StatusBacklog, model.StatusBreak}).OrderBy("release_date", firestore.Asc)
	case "playing":
		query = query.Where("status", "==", model.StatusPlaying).OrderBy("updated_at", firestore.Desc)
	case "history":
		query = query.Where("status", "in", []model.GameStatus{model.StatusDone, model.StatusAbandoned, model.StatusWontPlay}).OrderBy("date_played", firestore.Desc)
	default:
		// Return all games if no specific view is requested, or handle error
		log.Printf("Unknown view: %s, returning all games for user", view)
	}

	iter := query.Documents(ctx)
	var games []model.Game
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
			return nil, fmt.Errorf("failed to decode game data: %w", err)
		}
		games = append(games, game)
	}
	return games, nil
}

func (c *Client) GetGame(ctx context.Context, gameID string) (*model.Game, error) {
	doc, err := c.client.Collection("games").Doc(gameID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	var game model.Game
	if err := doc.DataTo(&game); err != nil {
		return nil, fmt.Errorf("failed to decode game data: %w", err)
	}
	return &game, nil
}

func (c *Client) GetAllGames(ctx context.Context) ([]model.Game, error) {
	iter := c.client.Collection("games").Documents(ctx)
	var games []model.Game
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
			return nil, fmt.Errorf("failed to decode game data: %w", err)
		}
		games = append(games, game)
	}
	return games, nil
}

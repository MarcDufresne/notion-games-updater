package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/marc/notion-games-updater/internal/domain"
	"github.com/marc/notion-games-updater/internal/igdb"
	notionmapper "github.com/marc/notion-games-updater/internal/notion"
)

const (
	gameFieldsQuery = `fields name,url,aggregated_rating,category,first_release_date,` +
		`platforms.*,cover.*,genres.*,websites.*,game_type.*,` +
		`release_dates.*,release_dates.status.*,release_dates.platform.*,` +
		`parent_game.id,parent_game.name,url,updated_at;`

	searchFieldsQuery = `fields game.name,game.url,game.aggregated_rating,game.category,game.first_release_date,` +
		`game.platforms.*,game.cover.*,game.genres.*,game.websites.*,game.game_type.*,` +
		`game.release_dates.*,game.release_dates.status.*,game.release_dates.platform.*,` +
		`game.url,game.updated_at;`
)

type UpdateGamesUseCase struct {
	igdbClient   *igdb.Client
	notionClient *notionapi.Client
	databaseID   string
	dryRun       bool
}

func NewUpdateGamesUseCase(igdbClient *igdb.Client, notionClient *notionapi.Client, databaseID string, dryRun bool) *UpdateGamesUseCase {
	return &UpdateGamesUseCase{
		igdbClient:   igdbClient,
		notionClient: notionClient,
		databaseID:   databaseID,
		dryRun:       dryRun,
	}
}

// getTitleFromPage extracts the title from a Notion page, trying common property names
func getTitleFromPage(page *notionapi.Page) (string, error) {
	// Try "Game" property first (both pointer and non-pointer types)
	if prop, ok := page.Properties["Game"]; ok {
		if titleProp, ok := prop.(*notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
			return titleProp.Title[0].PlainText, nil
		}
		if titleProp, ok := prop.(notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
			return titleProp.Title[0].PlainText, nil
		}
	}

	// Try "Name" property (both pointer and non-pointer types)
	if prop, ok := page.Properties["Name"]; ok {
		if titleProp, ok := prop.(*notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
			return titleProp.Title[0].PlainText, nil
		}
		if titleProp, ok := prop.(notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
			return titleProp.Title[0].PlainText, nil
		}
	}

	// Try to find any title property (both pointer and non-pointer types)
	for _, prop := range page.Properties {
		if titleProp, ok := prop.(*notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
			return titleProp.Title[0].PlainText, nil
		}
		if titleProp, ok := prop.(notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
			return titleProp.Title[0].PlainText, nil
		}
	}

	return "", fmt.Errorf("no title property found")
}

func (uc *UpdateGamesUseCase) getGame(ctx context.Context, page *notionapi.Page) (*domain.Game, notionapi.Properties, error) {
	// Get page title
	pageTitle, err := getTitleFromPage(page)
	if err != nil {
		return nil, nil, err
	}

	properties := notionapi.Properties{}

	// Try to get IGDB ID from the page
	var igdbID string
	if igdbIDProp, ok := page.Properties[string(domain.PropIGDBID)].(*notionapi.RichTextProperty); ok {
		if len(igdbIDProp.RichText) > 0 {
			igdbID = igdbIDProp.RichText[0].PlainText
		}
	}

	var game *domain.Game

	if igdbID != "" {
		// Clear other results if ID doesn't start with ":"
		if !strings.HasPrefix(igdbID, ":") {
			properties[string(domain.PropOtherResults)] = notionmapper.ClearOtherResults()
		} else {
			igdbID = strings.TrimPrefix(igdbID, ":")
		}

		// Check if context is cancelled
		if ctx.Err() != nil {
			return nil, properties, ctx.Err()
		}

		// Fetch by ID
		query := fmt.Sprintf("%s where id = %s;", gameFieldsQuery, igdbID)
		body, err := uc.igdbClient.Request("games", query)
		if err != nil {
			return nil, properties, fmt.Errorf("failed to fetch game: %w", err)
		}

		var games []*domain.Game
		if err := json.Unmarshal(body, &games); err != nil {
			return nil, properties, fmt.Errorf("failed to unmarshal game response: %w", err)
		}

		if len(games) == 0 {
			log.Printf("WARNING: Game %s not found in IGDB, skipping page '%s'", igdbID, pageTitle)
			return nil, properties, nil
		}

		game = games[0]
	} else {
		// Check if context is cancelled
		if ctx.Err() != nil {
			return nil, properties, ctx.Err()
		}

		// Search by title
		query := fmt.Sprintf(`%s search "%s"; where game != null & game.game_type.type != (13) & game.version_parent = null;`,
			searchFieldsQuery, pageTitle)

		body, err := uc.igdbClient.Request("search", query)
		if err != nil {
			return nil, properties, fmt.Errorf("failed to search game: %w", err)
		}

		var results []*domain.SearchResult
		if err := json.Unmarshal(body, &results); err != nil {
			return nil, properties, fmt.Errorf("failed to unmarshal search response: %w", err)
		}

		if len(results) == 0 {
			log.Printf("WARNING: No results found on IGDB for page '%s', skipping page", pageTitle)
			return nil, properties, nil
		}

		if len(results) > 1 {
			log.Printf("WARNING: Multiple results found for %s; defaulting to the first one", pageTitle)

			// Extract games for the other results property
			var games []*domain.Game
			for _, res := range results {
				if res.Game != nil {
					games = append(games, res.Game)
				}
			}
			properties[string(domain.PropOtherResults)] = notionmapper.CreateOtherResultsProperty(games)
		}

		game = results[0].Game
		if game == nil {
			log.Printf("WARNING: Search result has no game for page '%s', skipping page", pageTitle)
			return nil, properties, nil
		}
	}

	return game, properties, nil
}

func (uc *UpdateGamesUseCase) updatePage(ctx context.Context, pageID string) error {
	log.Printf("Processing page: %s", pageID)

	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Retrieve the page
	page, err := uc.notionClient.Page.Get(ctx, notionapi.PageID(pageID))
	if err != nil {
		return fmt.Errorf("failed to retrieve page: %w", err)
	}

	// Get page title
	pageTitle, err := getTitleFromPage(page)
	if err != nil {
		// Debug: log all available properties
		log.Printf("DEBUG: Available properties for page %s:", pageID)
		for propName, prop := range page.Properties {
			log.Printf("  - %s: %T", propName, prop)
		}
		return fmt.Errorf("failed to get page title: %w", err)
	}

	log.Printf("Game: %s", pageTitle)

	// Get the game data
	game, initialProps, err := uc.getGame(ctx, page)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}
	if game == nil {
		log.Printf("WARNING: Game not found for page '%s'", pageTitle)
		return nil
	}

	// Check if context is cancelled before updating
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Map game to Notion properties
	gameProps := notionmapper.MapGameToNotionProperties(game)

	// Merge initial properties with game properties
	for key, val := range initialProps {
		gameProps[key] = val
	}

	// Update the page
	updateReq := &notionapi.PageUpdateRequest{
		Properties: gameProps,
	}

	// Add cover if available
	if game.Cover != nil {
		coverURL := game.Cover.CoverBig2xURL()
		updateReq.Cover = &notionapi.Image{
			Type: notionapi.FileTypeExternal,
			External: &notionapi.FileObject{
				URL: coverURL,
			},
		}
	}

	// Skip update if in dry-run mode
	if uc.dryRun {
		log.Printf("DRY-RUN: Would update page '%s' (ID: %s)", pageTitle, pageID)
		if game.Cover != nil {
			log.Printf("DRY-RUN: Would set cover to: %s", game.Cover.CoverBig2xURL())
		}
		return nil
	}

	_, err = uc.notionClient.Page.Update(ctx, notionapi.PageID(pageID), updateReq)
	if err != nil {
		return fmt.Errorf("failed to update page: %w", err)
	}

	log.Printf("SUCCESS: Updated page '%s'", pageTitle)
	return nil
}

func (uc *UpdateGamesUseCase) Run(ctx context.Context) error {
	log.Println("Querying Notion database...")

	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Query the database
	resp, err := uc.notionClient.Database.Query(ctx, notionapi.DatabaseID(uc.databaseID), nil)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}

	log.Printf("Found %d pages to process", len(resp.Results))

	// Process each page
	for _, page := range resp.Results {
		// Check if context is cancelled before processing each page
		if ctx.Err() != nil {
			log.Println("Context cancelled, stopping page processing...")
			return ctx.Err()
		}

		if err := uc.updatePage(ctx, string(page.ID)); err != nil {
			// If error is due to context cancellation, return immediately
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("ERROR: Failed to update page %s: %v", page.ID, err)
			// Continue processing other pages
		}
	}

	log.Println("SUCCESS: Database updated")
	return nil
}

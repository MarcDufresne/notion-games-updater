package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"

	"game-tracker/internal/cache"
	"game-tracker/internal/database"
	"game-tracker/internal/igdb"
	"game-tracker/internal/legacy_domain"
	"game-tracker/internal/middleware"
	"game-tracker/internal/model"
)

type Handler struct {
	db         *database.Client
	igdbClient *igdb.Client
	cache      *cache.Cache
	authClient *auth.Client
}

func NewHandler(db *database.Client, igdbClient *igdb.Client, searchCache *cache.Cache, authClient *auth.Client) *Handler {
	return &Handler{
		db:         db,
		igdbClient: igdbClient,
		cache:      searchCache,
		authClient: authClient,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Create auth middleware
	authMW := middleware.AuthMiddleware(h.authClient)

	// Protected routes
	mux.Handle("/api/v1/games", authMW(http.HandlerFunc(h.handleGames)))
	mux.Handle("/api/v1/games/", authMW(http.HandlerFunc(h.handleGameByID)))
	mux.Handle("/api/v1/search", authMW(http.HandlerFunc(h.handleSearch)))
}

// handleGames handles GET /api/v1/games?view={backlog|playing|history}
func (h *Handler) handleGames(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getGames(w, r)
	} else if r.Method == http.MethodPost {
		h.createGame(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getGames retrieves games based on view parameter
func (h *Handler) getGames(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	view := r.URL.Query().Get("view")

	var games []*model.Game
	var err error

	switch view {
	case "backlog":
		games, err = h.db.GetBacklog(r.Context(), userID)
	case "playing":
		games, err = h.db.GetPlaying(r.Context(), userID)
	case "history":
		games, err = h.db.GetHistory(r.Context(), userID)
	case "calendar":
		games, err = h.db.GetUpcoming(r.Context(), userID)
	case "":
		// Return all games if no view specified
		games, err = h.db.GetGames(r.Context(), userID)
	default:
		http.Error(w, "Invalid view parameter", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("ERROR: Failed to fetch games: %v", err)
		http.Error(w, "Failed to fetch games", http.StatusInternalServerError)
		return
	}

	respondJSON(w, games)
}

// CreateGameRequest represents a request to create a new game
type CreateGameRequest struct {
	Title     string           `json:"title"`
	IGDBID    *int             `json:"igdb_id,omitempty"`
	Status    model.GameStatus `json:"status,omitempty"`
	CoverURL  string           `json:"cover_url,omitempty"`
	Rating    int              `json:"rating,omitempty"`
	Genres    []string         `json:"genres,omitempty"`
	Platforms []string         `json:"platforms,omitempty"`
}

// createGame handles POST /api/v1/games
func (h *Handler) createGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	game := &model.Game{
		UserID: userID,
		Title:  req.Title,
		IGDBID: req.IGDBID,
		Status: req.Status,
	}

	// Default status to Backlog if not provided
	if game.Status == "" {
		game.Status = model.StatusBacklog
	}

	// Validate status
	if !game.Status.IsValid() {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// If IGDB ID is provided, fetch metadata immediately
	if req.IGDBID != nil {
		igdbGame, err := h.igdbClient.GetGameByID(*req.IGDBID)
		if err != nil {
			log.Printf("ERROR: Failed to fetch IGDB metadata for game ID %d: %v", *req.IGDBID, err)
			http.Error(w, "Failed to fetch game metadata from IGDB", http.StatusInternalServerError)
			return
		}

		enrichGameFromIGDB(game, igdbGame)
	} else {
		// Use provided metadata for manual entries
		game.CoverURL = req.CoverURL
		game.Rating = req.Rating
		game.Genres = req.Genres
		game.Platforms = req.Platforms
	}

	// Save game
	if err := h.db.SaveGame(r.Context(), game); err != nil {
		log.Printf("ERROR: Failed to save game: %v", err)
		http.Error(w, "Failed to save game", http.StatusInternalServerError)
		return
	}

	respondJSON(w, game)
}

// handleGameByID handles routes like /api/v1/games/{id}/status
func (h *Handler) handleGameByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse URL path: /api/v1/games/{id}/status or /api/v1/games/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/games/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Game ID required", http.StatusBadRequest)
		return
	}

	gameID := parts[0]

	// Check if this is a status update request
	if len(parts) == 2 && parts[1] == "status" && r.Method == http.MethodPost {
		h.updateGameStatus(w, r, userID, gameID)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// UpdateStatusRequest represents a request to update game status
type UpdateStatusRequest struct {
	Status model.GameStatus `json:"status"`
}

// updateGameStatus handles POST /api/v1/games/{id}/status
func (h *Handler) updateGameStatus(w http.ResponseWriter, r *http.Request, userID, gameID string) {
	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.Status.IsValid() {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// Verify game belongs to user
	game, err := h.db.GetGame(r.Context(), gameID)
	if err != nil {
		log.Printf("ERROR: Failed to fetch game: %v", err)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if game.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Update status
	if err := h.db.UpdateGameStatus(r.Context(), gameID, req.Status); err != nil {
		log.Printf("ERROR: Failed to update game status: %v", err)
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	// Fetch updated game
	game, err = h.db.GetGame(r.Context(), gameID)
	if err != nil {
		log.Printf("ERROR: Failed to fetch updated game: %v", err)
		http.Error(w, "Failed to fetch updated game", http.StatusInternalServerError)
		return
	}

	respondJSON(w, game)
}

// handleSearch handles GET /api/v1/search?q={query}
func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Check cache first
	if results, found := h.cache.Get(query); found {
		log.Printf("Cache hit for query: %s", query)
		respondJSON(w, results)
		return
	}

	// Cache miss, query IGDB
	log.Printf("Cache miss for query: %s", query)
	results, err := h.igdbClient.Search(query)
	if err != nil {
		log.Printf("ERROR: Failed to search IGDB: %v", err)
		http.Error(w, "Failed to search games", http.StatusInternalServerError)
		return
	}

	// Store in cache
	h.cache.Set(query, results)

	respondJSON(w, results)
}

// enrichGameFromIGDB populates game fields from IGDB data
func enrichGameFromIGDB(game *model.Game, igdbGame *legacy_domain.Game) {
	// Update title if not manually set
	if game.Title == "" {
		game.Title = igdbGame.Name
	}

	// Set cover URL
	if igdbGame.Cover != nil {
		game.CoverURL = igdbGame.Cover.CoverBig2xURL()
	}

	// Set rating
	if igdbGame.AggregatedRating != nil {
		game.Rating = int(*igdbGame.AggregatedRating)
	}

	// Set genres
	if len(igdbGame.Genres) > 0 {
		game.Genres = make([]string, len(igdbGame.Genres))
		for i, genre := range igdbGame.Genres {
			game.Genres[i] = genre.Name
		}
	}

	// Set platforms
	if len(igdbGame.Platforms) > 0 {
		game.Platforms = make([]string, 0, len(igdbGame.Platforms))
		for _, platform := range igdbGame.Platforms {
			// Skip Stadia
			if platform.Abbreviation == "Stadia" {
				continue
			}
			game.Platforms = append(game.Platforms, platform.Abbreviation)
		}
	}

	// Set release date
	if igdbGame.FirstReleaseDate != nil {
		releaseDate := time.Unix(*igdbGame.FirstReleaseDate, 0)
		game.ReleaseDate = &releaseDate
	}
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ERROR: Failed to encode JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

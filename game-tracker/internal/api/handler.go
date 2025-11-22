package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"game-tracker/internal/database"
	"game-tracker/internal/middleware"
	"game-tracker/internal/model"
)

type Handler struct {
	db *database.Client
}

// NewHandler creates a new API handler
func NewHandler(db *database.Client) *Handler {
	return &Handler{
		db: db,
	}
}

// GetGames handles GET /api/v1/games?view={backlog|playing|history}
func (h *Handler) GetGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
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
	case "":
		// No filter - get all games
		games, err = h.db.GetGames(r.Context(), userID, nil)
	default:
		http.Error(w, fmt.Sprintf("Invalid view: %s", view), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get games: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// CreateOrUpdateGame handles POST /api/v1/games
func (h *Handler) CreateOrUpdateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var game model.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Ensure user_id is set correctly (will be overwritten by SaveGame anyway)
	game.UserID = userID

	if err := h.db.SaveGame(r.Context(), &game); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save game: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

// UpdateGameStatus handles POST /api/v1/games/{id}/status
func (h *Handler) UpdateGameStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract game ID from path
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	gameID := parts[3] // /api/v1/games/{id}/status

	// Parse request body
	var req struct {
		Status     model.GameStatus `json:"status"`
		DatePlayed *time.Time       `json:"date_played,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Get the existing game
	game, err := h.db.GetGame(r.Context(), userID, gameID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Game not found: %v", err), http.StatusNotFound)
		return
	}

	// Update status
	game.Status = req.Status

	// If moving to completed status, set date_played if provided
	if game.IsComplete() && req.DatePlayed != nil {
		game.DatePlayed = req.DatePlayed
	} else if game.IsComplete() && game.DatePlayed == nil {
		// Set to now if not provided
		now := time.Now()
		game.DatePlayed = &now
	}

	// Save the updated game
	if err := h.db.SaveGame(r.Context(), game); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

// HealthCheck handles GET /api/v1/health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

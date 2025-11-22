package api

import (
	"encoding/json"
	"net/http"

	"game-tracker/internal/database"
	"game-tracker/internal/middleware"
	"game-tracker/internal/model"
)

type Handler struct {
	db *database.Client
}

func NewHandler(db *database.Client) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/games", h.handleGetGames)
	mux.HandleFunc("POST /api/v1/games", h.handleSaveGame)
	mux.HandleFunc("POST /api/v1/games/{id}/status", h.handleUpdateStatus)
}

func (h *Handler) handleGetGames(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	view := r.URL.Query().Get("view")
	games, err := h.db.GetGames(r.Context(), userID, view)
	if err != nil {
		http.Error(w, "Failed to fetch games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func (h *Handler) handleSaveGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var game model.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	game.UserID = userID // Enforce user ID from token

	if err := h.db.SaveGame(r.Context(), &game); err != nil {
		http.Error(w, "Failed to save game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

func (h *Handler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKeyUID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		Status model.GameStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ideally, we should verify the game belongs to the user before updating.
	// For now, we'll just fetch it, update status, and save.
	// Wait, GetGames returns a list. I might need GetGame(id) or just rely on SaveGame overwrite.
	// But SaveGame overwrites everything. I need to read-modify-write.
	// For now, to proceed quickly, I'll assume the client sends the full object or I'll implement GetGame.
	// Or better, I'll implement a partial update if Firestore supports it easily, or just read-modify-write.

	game, err := h.db.GetGame(r.Context(), id)
	if err != nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if game.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	game.Status = payload.Status
	if err := h.db.SaveGame(r.Context(), game); err != nil {
		http.Error(w, "Failed to update game status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

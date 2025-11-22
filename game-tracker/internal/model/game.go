package model

import (
	"time"
)

// GameStatus represents the status of a game in the user's library
type GameStatus string

const (
	StatusBacklog   GameStatus = "Backlog"
	StatusBreak     GameStatus = "Break"      // Group: To-do
	StatusPlaying   GameStatus = "Playing"    // Group: In Progress
	StatusDone      GameStatus = "Done"       // Group: Complete
	StatusAbandoned GameStatus = "Abandoned"  // Group: Complete
	StatusWontPlay  GameStatus = "Won't Play" // Group: Complete
)

// Game represents a game in the user's library
type Game struct {
	ID          string     `firestore:"id" json:"id"`
	UserID      string     `firestore:"user_id" json:"user_id"`
	Title       string     `firestore:"title" json:"title"`
	IGDBID      *string    `firestore:"igdb_id" json:"igdb_id"` // Surrogate Key. Null = Manual Entry
	CoverURL    string     `firestore:"cover_url" json:"cover_url"`
	Rating      int        `firestore:"rating" json:"rating"`                           // 0-100
	Status      GameStatus `firestore:"status" json:"status"`
	Genres      []string   `firestore:"genres" json:"genres"`
	Platforms   []string   `firestore:"platforms" json:"platforms"`
	ReleaseDate *time.Time `firestore:"release_date" json:"release_date"`
	DatePlayed  *time.Time `firestore:"date_played" json:"date_played"`
	CreatedAt   time.Time  `firestore:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `firestore:"updated_at" json:"updated_at"`
}

// IsComplete returns true if the game is in a completed state
func (g *Game) IsComplete() bool {
	return g.Status == StatusDone || g.Status == StatusAbandoned || g.Status == StatusWontPlay
}

// IsInProgress returns true if the game is currently being played
func (g *Game) IsInProgress() bool {
	return g.Status == StatusPlaying
}

// IsTodo returns true if the game is in the backlog or on break
func (g *Game) IsTodo() bool {
	return g.Status == StatusBacklog || g.Status == StatusBreak
}

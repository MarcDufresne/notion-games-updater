package model

import (
	"time"
)

type GameStatus string

const (
	StatusBacklog   GameStatus = "Backlog"
	StatusBreak     GameStatus = "Break"      // Group: To-do
	StatusPlaying   GameStatus = "Playing"    // Group: In Progress
	StatusDone      GameStatus = "Done"       // Group: Complete
	StatusAbandoned GameStatus = "Abandoned"  // Group: Complete
	StatusWontPlay  GameStatus = "Won't Play" // Group: Complete
)

type Game struct {
	ID          string     `json:"id" firestore:"id"`
	UserID      string     `json:"user_id" firestore:"user_id"` // Partition Key
	Title       string     `json:"title" firestore:"title"`
	IGDBID      *string    `json:"igdb_id" firestore:"igdb_id"` // Surrogate Key. Null = Manual Entry.
	CoverURL    string     `json:"cover_url" firestore:"cover_url"`
	Rating      int        `json:"rating" firestore:"rating"` // Int (0-100)
	Status      GameStatus `json:"status" firestore:"status"`
	Genres      []string   `json:"genres" firestore:"genres"`
	Platforms   []string   `json:"platforms" firestore:"platforms"`
	ReleaseDate *time.Time `json:"release_date" firestore:"release_date"`
	DatePlayed  *time.Time `json:"date_played" firestore:"date_played"`
	CreatedAt   time.Time  `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" firestore:"updated_at"`
}

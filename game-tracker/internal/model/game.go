package model

import (
	"time"
)

type GameStatus string

const (
	StatusBacklog   GameStatus = "Backlog"
	StatusBreak     GameStatus = "Break"
	StatusPlaying   GameStatus = "Playing"
	StatusDone      GameStatus = "Done"
	StatusAbandoned GameStatus = "Abandoned"
	StatusWontPlay  GameStatus = "Won't Play"
)

// IsValid checks if the status is one of the valid values
func (s GameStatus) IsValid() bool {
	switch s {
	case StatusBacklog, StatusBreak, StatusPlaying, StatusDone, StatusAbandoned, StatusWontPlay:
		return true
	default:
		return false
	}
}

type Game struct {
	ID            string     `firestore:"id" json:"id"`
	UserID        string     `firestore:"user_id" json:"user_id"`
	Title         string     `firestore:"title" json:"title"`
	IGDBID        *int       `firestore:"igdb_id,omitempty" json:"igdb_id,omitempty"`
	CoverURL      string     `firestore:"cover_url,omitempty" json:"cover_url,omitempty"`
	Rating        int        `firestore:"rating,omitempty" json:"rating,omitempty"` // 0-100
	Status        GameStatus `firestore:"status" json:"status"`
	Genres        []string   `firestore:"genres,omitempty" json:"genres,omitempty"`
	Platforms     []string   `firestore:"platforms,omitempty" json:"platforms,omitempty"`
	ReleaseDate   *time.Time `firestore:"release_date,omitempty" json:"release_date,omitempty"`
	DatePlayed    *time.Time `firestore:"date_played,omitempty" json:"date_played,omitempty"`
	CreatedAt     time.Time  `firestore:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `firestore:"updated_at" json:"updated_at"`
	LastSyncError string     `firestore:"last_sync_error,omitempty" json:"last_sync_error,omitempty"`
}

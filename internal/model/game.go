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

type MatchStatus string

const (
	MatchStatusMatched     MatchStatus = "matched"      // Game has valid IGDB ID
	MatchStatusUnmatched   MatchStatus = "unmatched"    // Manual entry without IGDB ID
	MatchStatusMultiple    MatchStatus = "multiple"     // Multiple potential matches found
	MatchStatusNoMatch     MatchStatus = "no_match"     // Search found no matches
	MatchStatusNeedsReview MatchStatus = "needs_review" // User needs to review/fix match
)

type Game struct {
	ID            string      `firestore:"id" json:"id"`
	UserID        string      `firestore:"user_id" json:"user_id"`
	Title         string      `firestore:"title" json:"title"`
	IGDBID        int         `firestore:"igdb_id" json:"igdb_id"` // 0 means no IGDB ID (unmatched)
	CoverURL      string      `firestore:"cover_url,omitempty" json:"cover_url,omitempty"`
	Rating        int         `firestore:"rating,omitempty" json:"rating,omitempty"` // 0-100
	Status        GameStatus  `firestore:"status" json:"status"`
	Genres        []string    `firestore:"genres,omitempty" json:"genres,omitempty"`
	Platforms     []string    `firestore:"platforms,omitempty" json:"platforms,omitempty"`
	ReleaseDate   *time.Time  `firestore:"release_date,omitempty" json:"release_date,omitempty"`
	DatePlayed    *time.Time  `firestore:"date_played,omitempty" json:"date_played,omitempty"`
	SteamURL      string      `firestore:"steam_url,omitempty" json:"steam_url,omitempty"`
	OfficialURL   string      `firestore:"official_url,omitempty" json:"official_url,omitempty"`
	MatchStatus   MatchStatus `firestore:"match_status,omitempty" json:"match_status,omitempty"`
	CreatedAt     time.Time   `firestore:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `firestore:"updated_at" json:"updated_at"`
	LastSyncError string      `firestore:"last_sync_error,omitempty" json:"last_sync_error,omitempty"`
}

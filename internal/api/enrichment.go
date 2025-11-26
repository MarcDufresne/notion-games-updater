package api

import (
	"time"

	"game-tracker/internal/igdb"
	"game-tracker/internal/model"
)

func EnrichGameFromIGDB(game *model.Game, igdbGame *igdb.Game, trackChanges bool) bool {
	changed := false

	if newValue := igdbGame.Name; !trackChanges || game.Title != newValue {
		game.Title = newValue
		changed = true
	}

	if igdbGame.Cover != nil {
		if newValue := igdbGame.Cover.CoverBig2xURL(); !trackChanges || game.CoverURL != newValue {
			game.CoverURL = newValue
			changed = true
		}
	}

	if igdbGame.AggregatedRating != nil {
		if newValue := int(*igdbGame.AggregatedRating); !trackChanges || game.Rating != newValue {
			game.Rating = newValue
			changed = true
		}
	}

	if len(igdbGame.Genres) > 0 {
		newGenres := make([]string, len(igdbGame.Genres))
		for i, genre := range igdbGame.Genres {
			newGenres[i] = genre.Name
		}
		if !trackChanges || !stringSlicesEqual(game.Genres, newGenres) {
			game.Genres = newGenres
			changed = true
		}
	}

	if len(igdbGame.Platforms) > 0 {
		newPlatforms := make([]string, 0, len(igdbGame.Platforms))
		for _, platform := range igdbGame.Platforms {
			if platform.Abbreviation != "" && platform.Abbreviation != "Stadia" {
				newPlatforms = append(newPlatforms, platform.Abbreviation)
			}
		}
		if !trackChanges || !stringSlicesEqual(game.Platforms, newPlatforms) {
			game.Platforms = newPlatforms
			changed = true
		}
	}

	if igdbGame.FirstReleaseDate != nil {
		releaseDate := time.Unix(*igdbGame.FirstReleaseDate, 0)
		if !trackChanges || game.ReleaseDate == nil || !game.ReleaseDate.Equal(releaseDate) {
			game.ReleaseDate = &releaseDate
			changed = true
		}
	}

	if len(igdbGame.Websites) > 0 {
		newSteamURL := ""
		newOfficialURL := ""
		for _, website := range igdbGame.Websites {
			switch website.Type {
			case igdb.WebsiteCategorySteam:
				newSteamURL = website.URL
			case igdb.WebsiteCategoryOfficial:
				newOfficialURL = website.URL
			}
		}
		if !trackChanges || game.SteamURL != newSteamURL {
			game.SteamURL = newSteamURL
			changed = true
		}
		if !trackChanges || game.OfficialURL != newOfficialURL {
			game.OfficialURL = newOfficialURL
			changed = true
		}
	}

	if trackChanges && game.LastSyncError != "" {
		game.LastSyncError = ""
		changed = true
	}

	return changed
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

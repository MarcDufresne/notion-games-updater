package notion

import (
	"fmt"
	"sort"
	"time"

	"github.com/jomei/notionapi"
	"github.com/marc/notion-games-updater/internal/domain"
)

const notionTimezoneName = "America/New_York"

var releaseDateStatusMapping = map[domain.ReleaseDateStatusType]int{
	domain.ReleaseStatusFullRelease: 1,
	domain.ReleaseStatusEarlyAccess: 2,
}

// MapGameToNotionProperties converts a Game to Notion page properties
func MapGameToNotionProperties(game *domain.Game) notionapi.Properties {
	props := notionapi.Properties{
		string(domain.PropIGDBID): notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: fmt.Sprintf(":%d", game.ID),
					},
				},
			},
		},
		"Game": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: game.Name,
					},
				},
			},
		},
	}

	// Rating
	if game.AggregatedRating != nil {
		rating := *game.AggregatedRating / 100.0
		// Round to 2 decimal places
		rating = float64(int(rating*100+0.5)) / 100.0
		props[string(domain.PropRating)] = notionapi.NumberProperty{
			Number: rating,
		}
	}

	// Platforms
	if len(game.Platforms) > 0 {
		var platforms []notionapi.Option
		for _, platform := range game.Platforms {
			// Skip Stadia
			if platform.Abbreviation == "Stadia" {
				continue
			}
			platforms = append(platforms, notionapi.Option{
				Name: platform.Abbreviation,
			})
		}
		if len(platforms) > 0 {
			props[string(domain.PropPlatforms)] = notionapi.MultiSelectProperty{
				MultiSelect: platforms,
			}
		}
	}

	// Genres
	if len(game.Genres) > 0 {
		var genres []notionapi.Option
		for _, genre := range game.Genres {
			genres = append(genres, notionapi.Option{
				Name: genre.Name,
			})
		}
		props[string(domain.PropGenres)] = notionapi.MultiSelectProperty{
			MultiSelect: genres,
		}
	}

	// Store URL
	if len(game.Websites) > 0 {
		priorityOrder := []domain.WebsiteCategory{
			domain.WebsiteCategorySteam,
			domain.WebsiteCategoryItch,
			domain.WebsiteCategoryGOG,
			domain.WebsiteCategoryEpicGames,
		}
		categoryToURL := make(map[domain.WebsiteCategory]string)
		for _, website := range game.Websites {
			categoryToURL[website.Category] = website.URL
		}

		var url string
		for _, category := range priorityOrder {
			if u, ok := categoryToURL[category]; ok {
				url = u
				break
			}
		}

		if url != "" {
			props[string(domain.PropStoreURL)] = notionapi.URLProperty{
				URL: url,
			}
		}
	}

	// Release Date
	if len(game.ReleaseDates) > 0 {
		notionTimezone, err := time.LoadLocation(notionTimezoneName)
		if err != nil {
			notionTimezone = time.UTC
		}

		type dateMapping struct {
			priority int
			date     time.Time
			human    string
		}

		var dates []dateMapping

		for _, release := range game.ReleaseDates {
			// Filter by region
			if release.ReleaseRegion != domain.ReleaseRegionWorldwide && release.ReleaseRegion != domain.ReleaseRegionNorthAmerica {
				continue
			}

			// Filter by status
			if release.Status != nil &&
				release.Status.Name != domain.ReleaseStatusFullRelease &&
				release.Status.Name != domain.ReleaseStatusEarlyAccess {
				continue
			}

			// Filter out Stadia
			if release.Platform != nil && release.Platform.Abbreviation == "Stadia" {
				continue
			}

			// Only process categories 0-6
			if release.DateFormat < 0 || release.DateFormat > 6 {
				continue
			}

			var releaseDate time.Time

			switch release.DateFormat {
			case domain.ReleaseCategoryYYYYMMMMDD:
				if release.Date != nil {
					unix := time.Unix(*release.Date, 0).In(time.UTC)
					releaseDate = time.Date(unix.Year(), unix.Month(), unix.Day(), 0, 0, 0, 0, notionTimezone)
				}
			case domain.ReleaseCategoryYYYYMMMM:
				if release.Y != nil && release.M != nil {
					// End of month
					releaseDate = time.Date(*release.Y, time.Month(*release.M+1), 0, 0, 0, 0, 0, notionTimezone)
				}
			case domain.ReleaseCategoryYYYY:
				if release.Y != nil {
					// End of year
					releaseDate = time.Date(*release.Y, 12, 31, 0, 0, 0, 0, notionTimezone)
				}
			case domain.ReleaseCategoryYYYYQ1:
				if release.Y != nil {
					// End of Q1 (March)
					releaseDate = time.Date(*release.Y, 3, 31, 0, 0, 0, 0, notionTimezone)
				}
			case domain.ReleaseCategoryYYYYQ2:
				if release.Y != nil {
					// End of Q2 (June)
					releaseDate = time.Date(*release.Y, 6, 30, 0, 0, 0, 0, notionTimezone)
				}
			case domain.ReleaseCategoryYYYYQ3:
				if release.Y != nil {
					// End of Q3 (September)
					releaseDate = time.Date(*release.Y, 9, 30, 0, 0, 0, 0, notionTimezone)
				}
			case domain.ReleaseCategoryYYYYQ4:
				if release.Y != nil {
					// End of Q4 (December)
					releaseDate = time.Date(*release.Y, 12, 31, 0, 0, 0, 0, notionTimezone)
				}
			}

			if releaseDate.IsZero() {
				continue
			}

			priority := 0
			if release.Status != nil {
				if p, ok := releaseDateStatusMapping[release.Status.Name]; ok {
					priority = p
				}
			}

			dates = append(dates, dateMapping{
				priority: priority,
				date:     releaseDate,
				human:    release.Human,
			})
		}

		if len(dates) > 0 {
			// Sort by priority first, then by date
			sort.Slice(dates, func(i, j int) bool {
				if dates[i].priority != dates[j].priority {
					return dates[i].priority < dates[j].priority
				}
				return dates[i].date.Before(dates[j].date)
			})

			selectedDate := dates[0]
			dateValue := notionapi.Date(selectedDate.date)

			props[string(domain.PropReleaseDate)] = notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: &dateValue,
				},
			}

			if selectedDate.human != "" {
				props[string(domain.PropReleaseDateHuman)] = notionapi.RichTextProperty{
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: selectedDate.human,
							},
						},
					},
				}
			}
		} else {
			// No release date found
			props[string(domain.PropReleaseDate)] = notionapi.DateProperty{
				Date: nil,
			}
			props[string(domain.PropReleaseDateHuman)] = notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: "TBD",
						},
					},
				},
			}
		}
	}

	return props
}

// ClearOtherResults clears the "Other Results" property
func ClearOtherResults() notionapi.Property {
	return notionapi.RichTextProperty{
		RichText: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: "",
				},
			},
		},
	}
}

// CreateOtherResultsProperty creates the "Other Results" property with multiple game results
func CreateOtherResultsProperty(games []*domain.Game) notionapi.Property {
	var richTexts []notionapi.RichText

	for _, game := range games {
		richTexts = append(richTexts,
			notionapi.RichText{
				Text: &notionapi.Text{
					Content: game.Name,
					Link: &notionapi.Link{
						Url: game.URL,
					},
				},
			},
			notionapi.RichText{
				Text: &notionapi.Text{
					Content: fmt.Sprintf(" (%d)\n", game.ID),
				},
			},
		)
	}

	return notionapi.RichTextProperty{
		RichText: richTexts,
	}
}

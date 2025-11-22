package domain

type Cover struct {
	ID      int    `json:"id"`
	ImageID string `json:"image_id"`
}

func (c *Cover) CoverBig2xURL() string {
	return "https://images.igdb.com/igdb/image/upload/t_cover_big_2x/" + c.ImageID + ".jpg"
}

type Platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Website struct {
	Category WebsiteCategory `json:"category"`
	URL      string          `json:"url"`
}

type ReleaseDateStatus struct {
	ID   int                   `json:"id"`
	Name ReleaseDateStatusType `json:"name"`
}

type ReleaseDate struct {
	DateFormat    ReleaseCategory    `json:"date_format"`
	Date          *int64             `json:"date,omitempty"`
	Human         string             `json:"human,omitempty"`
	Y             *int               `json:"y,omitempty"`
	M             *int               `json:"m,omitempty"`
	ReleaseRegion ReleaseRegion      `json:"release_region"`
	Status        *ReleaseDateStatus `json:"status,omitempty"`
	Platform      *Platform          `json:"platform,omitempty"`
}

type GameType struct {
	Type string `json:"type"`
}

type Game struct {
	ID               int           `json:"id"`
	AggregatedRating *float64      `json:"aggregated_rating,omitempty"`
	Cover            *Cover        `json:"cover,omitempty"`
	FirstReleaseDate *int64        `json:"first_release_date,omitempty"`
	GameType         *GameType     `json:"game_type,omitempty"`
	Genres           []Genre       `json:"genres,omitempty"`
	Name             string        `json:"name"`
	Platforms        []Platform    `json:"platforms,omitempty"`
	ReleaseDates     []ReleaseDate `json:"release_dates,omitempty"`
	UpdatedAt        *int64        `json:"updated_at,omitempty"`
	URL              string        `json:"url,omitempty"`
	Websites         []Website     `json:"websites,omitempty"`
}

type SearchResult struct {
	Game *Game `json:"game,omitempty"`
}

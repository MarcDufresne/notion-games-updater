package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"game-tracker/internal/legacy_domain"
)

const (
	authURL = "https://id.twitch.tv/oauth2/token"
	apiURL  = "https://api.igdb.com/v4/"
	retries = 3
)

type Client struct {
	clientID           string
	clientSecret       string
	accessToken        string
	accessTokenExpires int64
	httpClient         *http.Client
	mu                 sync.Mutex
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// SearchCandidate represents a minimal game result for search
type SearchCandidate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CoverURL    string `json:"cover_url"`
	ReleaseYear int    `json:"release_year"`
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) refreshToken() error {
	params := fmt.Sprintf("?client_id=%s&client_secret=%s&grant_type=client_credentials",
		c.clientID, c.clientSecret)

	log.Println("[IGDB] Refreshing access token...")
	resp, err := c.httpClient.Post(authURL+params, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[IGDB] Token refresh failed with status %d: %s", resp.StatusCode, string(body))
		return fmt.Errorf("failed to refresh token [%d]: %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.accessTokenExpires = time.Now().Unix() + int64(tokenResp.ExpiresIn)

	log.Printf("[IGDB] Token refreshed successfully (expires in %d seconds)", tokenResp.ExpiresIn)
	return nil
}

func (c *Client) request(endpoint, query string) (*http.Response, error) {
	c.mu.Lock()
	if c.accessToken == "" || c.accessTokenExpires-3600 < time.Now().Unix() {
		if err := c.refreshToken(); err != nil {
			c.mu.Unlock()
			return nil, err
		}
	}
	token := c.accessToken
	c.mu.Unlock()

	req, err := http.NewRequest("POST", apiURL+endpoint, bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Client-ID", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "text/plain")

	var resp *http.Response
	tryCount := 0

	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			tryCount = 0
			log.Printf("[IGDB] Rate limited, retrying after 250ms...")
			time.Sleep(250 * time.Millisecond)
			continue
		}

		tryCount++
		if tryCount > retries {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("[IGDB] Request failed after %d retries [%d]: %s", retries, resp.StatusCode, string(body))
			return nil, fmt.Errorf("failed to make request [%d]: %s", resp.StatusCode, string(body))
		}
		resp.Body.Close()
		log.Printf("[IGDB] Request failed with status %d, retry %d/%d", resp.StatusCode, tryCount, retries)
	}

	return resp, nil
}

func (c *Client) Request(endpoint, query string) ([]byte, error) {
	resp, err := c.request(endpoint, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// Search performs a search query on IGDB and returns minimal candidate results
func (c *Client) Search(query string) ([]SearchCandidate, error) {
	log.Printf("[IGDB] Searching for: %q", query)
	searchQuery := fmt.Sprintf(`fields game.name,game.cover.*,game.first_release_date; search "%s"; where game != null & game.game_type.type != (13) & game.version_parent = null; limit 10;`, query)

	body, err := c.Request("search", searchQuery)
	if err != nil {
		log.Printf("[IGDB] Search failed for %q: %v", query, err)
		return nil, fmt.Errorf("failed to search games: %w", err)
	}

	var results []legacy_domain.SearchResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	candidates := make([]SearchCandidate, 0, len(results))
	for _, result := range results {
		if result.Game == nil {
			continue
		}

		candidate := SearchCandidate{
			ID:   result.Game.ID,
			Name: result.Game.Name,
		}

		if result.Game.Cover != nil {
			candidate.CoverURL = result.Game.Cover.CoverBig2xURL()
		}

		if result.Game.FirstReleaseDate != nil {
			releaseTime := time.Unix(*result.Game.FirstReleaseDate, 0)
			candidate.ReleaseYear = releaseTime.Year()
		}

		candidates = append(candidates, candidate)
	}

	log.Printf("[IGDB] Search for %q returned %d results", query, len(candidates))
	return candidates, nil
}

// GetGameByID fetches full game details by IGDB ID
func (c *Client) GetGameByID(id int) (*legacy_domain.Game, error) {
	log.Printf("[IGDB] Fetching game details for ID: %d", id)
	query := fmt.Sprintf(`fields name,url,aggregated_rating,category,first_release_date,platforms.*,cover.*,genres.*,websites.*,game_type.*,release_dates.*,release_dates.status.*,release_dates.platform.*,parent_game.id,parent_game.name,url,updated_at; where id = %d;`, id)

	body, err := c.Request("games", query)
	if err != nil {
		log.Printf("[IGDB] Failed to fetch game ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch game: %w", err)
	}

	var games []*legacy_domain.Game
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game response: %w", err)
	}

	if len(games) == 0 {
		log.Printf("[IGDB] Game ID %d not found", id)
		return nil, fmt.Errorf("game with ID %d not found", id)
	}

	log.Printf("[IGDB] Successfully fetched game: %s (ID: %d)", games[0].Name, id)

	return games[0], nil
}

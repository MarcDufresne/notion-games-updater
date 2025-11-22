package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
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

	resp, err := c.httpClient.Post(authURL+params, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to refresh token [%d]: %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.accessTokenExpires = time.Now().Unix() + int64(tokenResp.ExpiresIn)

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
			time.Sleep(250 * time.Millisecond)
			continue
		}

		tryCount++
		if tryCount > retries {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("failed to make request [%d]: %s", resp.StatusCode, string(body))
		}
		resp.Body.Close()
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

func (c *Client) ListAll(endpoint, query string) ([]byte, error) {
	resp, err := c.request(endpoint, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data []json.RawMessage
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("expected a list response: %w", err)
	}

	totalStr := resp.Header.Get("x-count")
	if totalStr == "" {
		return body, nil
	}

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return body, nil
	}

	offset := len(data)
	allData := data

	for offset < total {
		newQuery := fmt.Sprintf("%s offset %d;", query, offset)
		fmt.Println(newQuery)

		resp, err := c.request(endpoint, newQuery)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var newData []json.RawMessage
		if err := json.Unmarshal(body, &newData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal paginated response: %w", err)
		}

		allData = append(allData, newData...)
		offset += len(newData)
	}

	return json.Marshal(allData)
}

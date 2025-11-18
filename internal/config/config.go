package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	IGDB struct {
		ClientID     string
		ClientSecret string
	}
	Notion struct {
		Token      string
		DatabaseID string
	}
	DryRun bool
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.IGDB.ClientID = os.Getenv("IGDB_CLIENT_ID")
	cfg.IGDB.ClientSecret = os.Getenv("IGDB_CLIENT_SECRET")
	cfg.Notion.Token = os.Getenv("NOTION_TOKEN")
	cfg.Notion.DatabaseID = os.Getenv("NOTION_DATABASE_ID")
	cfg.DryRun = os.Getenv("DRY_RUN") == "true" || os.Getenv("DRY_RUN") == "1"

	// Validate required fields
	if cfg.IGDB.ClientID == "" {
		return nil, fmt.Errorf("IGDB_CLIENT_ID is required")
	}
	if cfg.IGDB.ClientSecret == "" {
		return nil, fmt.Errorf("IGDB_CLIENT_SECRET is required")
	}
	if cfg.Notion.Token == "" {
		return nil, fmt.Errorf("NOTION_TOKEN is required")
	}
	if cfg.Notion.DatabaseID == "" {
		return nil, fmt.Errorf("NOTION_DATABASE_ID is required")
	}

	return cfg, nil
}

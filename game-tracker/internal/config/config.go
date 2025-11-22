package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Firebase struct {
		ProjectID           string
		ServiceAccountJSON  string // Path to service account JSON file
		ServiceAccountKey   []byte // Raw service account key (for Docker)
	}
	IGDB struct {
		ClientID     string
		ClientSecret string
	}
	Server struct {
		Port string
		Host string
	}
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{}

	// Firebase configuration
	cfg.Firebase.ProjectID = os.Getenv("FIREBASE_PROJECT_ID")
	cfg.Firebase.ServiceAccountJSON = os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")
	if rawKey := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY"); rawKey != "" {
		cfg.Firebase.ServiceAccountKey = []byte(rawKey)
	}

	// IGDB configuration
	cfg.IGDB.ClientID = os.Getenv("IGDB_CLIENT_ID")
	cfg.IGDB.ClientSecret = os.Getenv("IGDB_CLIENT_SECRET")

	// Server configuration
	cfg.Server.Port = os.Getenv("PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	cfg.Server.Host = os.Getenv("HOST")
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}

	// Validate required fields
	if cfg.Firebase.ProjectID == "" {
		return nil, fmt.Errorf("FIREBASE_PROJECT_ID is required")
	}
	if cfg.Firebase.ServiceAccountJSON == "" && len(cfg.Firebase.ServiceAccountKey) == 0 {
		return nil, fmt.Errorf("either FIREBASE_SERVICE_ACCOUNT_JSON or FIREBASE_SERVICE_ACCOUNT_KEY is required")
	}
	if cfg.IGDB.ClientID == "" {
		return nil, fmt.Errorf("IGDB_CLIENT_ID is required")
	}
	if cfg.IGDB.ClientSecret == "" {
		return nil, fmt.Errorf("IGDB_CLIENT_SECRET is required")
	}

	return cfg, nil
}

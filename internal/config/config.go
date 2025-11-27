package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Firebase struct {
		ProjectID          string
		ServiceAccountKey  string // Raw JSON key
		ServiceAccountJSON string // File path
	}
	IGDB struct {
		ClientID     string
		ClientSecret string
	}
	Server struct {
		Port   string
		Host   string
		NoSync bool
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.Firebase.ServiceAccountKey = os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	cfg.Firebase.ServiceAccountJSON = os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")
	cfg.Firebase.ProjectID = os.Getenv("FIREBASE_PROJECT_ID")

	if cfg.Firebase.ServiceAccountKey == "" && cfg.Firebase.ServiceAccountJSON == "" {
		return nil, fmt.Errorf("either FIREBASE_SERVICE_ACCOUNT_KEY or FIREBASE_SERVICE_ACCOUNT_JSON is required")
	}

	cfg.IGDB.ClientID = os.Getenv("IGDB_CLIENT_ID")
	cfg.IGDB.ClientSecret = os.Getenv("IGDB_CLIENT_SECRET")

	if cfg.IGDB.ClientID == "" {
		return nil, fmt.Errorf("IGDB_CLIENT_ID is required")
	}
	if cfg.IGDB.ClientSecret == "" {
		return nil, fmt.Errorf("IGDB_CLIENT_SECRET is required")
	}

	cfg.Server.Port = os.Getenv("PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	cfg.Server.Host = os.Getenv("HOST")
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}

	noSync := os.Getenv("NO_SYNC")
	if noSync == "1" || noSync == "true" {
		cfg.Server.NoSync = true
	}

	return cfg, nil
}

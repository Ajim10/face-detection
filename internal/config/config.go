package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ProjectID string
	Location  string
	Model     string
	Bucket    string
	Object    string
	Filename  string
}

func Load() (*Config, error) {
	cfg := &Config{
		ProjectID: os.Getenv("PROJECT_ID"),
		Location:  os.Getenv("LOCATION"),
		Model:     os.Getenv("MODEL"),
		Bucket:    os.Getenv("BUCKET"),
		Object:    os.Getenv("OBJECT"),
		Filename:  os.Getenv("FILENAME"),
	}

	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("PROJECT_ID is not set")
	}
	if cfg.Location == "" {
		return nil, fmt.Errorf("LOCATION is not set")
	}
	if cfg.Model == "" {
		return nil, fmt.Errorf("MODEL is not set")
	}
	return cfg, nil
}

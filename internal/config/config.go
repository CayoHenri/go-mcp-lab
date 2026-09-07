package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	OpenAIKey   string
}

func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL não configurada")
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return Config{}, fmt.Errorf("OPENAI_API_KEY não configurada")
	}

	return Config{
		DatabaseURL: databaseURL,
		OpenAIKey:   openAIKey,
	}, nil
}

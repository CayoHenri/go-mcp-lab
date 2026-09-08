package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load()
}

type ServerConfig struct {
	DatabaseURL string
}

func LoadServer() (ServerConfig, error) {
	LoadEnv()

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return ServerConfig{}, fmt.Errorf("DATABASE_URL não configurada")
	}

	return ServerConfig{DatabaseURL: databaseURL}, nil
}

type ClientConfig struct {
	OpenAIAPIKey string
	AgentModel   string
}

func LoadClient() (ClientConfig, error) {
	LoadEnv()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return ClientConfig{}, fmt.Errorf("OPENAI_API_KEY não configurada")
	}

	agentModel := os.Getenv("AGENT_MODEL")
	if agentModel == "" {
		return ClientConfig{}, fmt.Errorf("AGENT_MODEL não configurada")
	}

	return ClientConfig{
		OpenAIAPIKey: apiKey,
		AgentModel:   agentModel,
	}, nil
}

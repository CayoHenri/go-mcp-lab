package main

import (
	"context"
	"fmt"
	"log"

	"github.com/CayoHenri/go-mcp-lab/internal/agent"
	"github.com/CayoHenri/go-mcp-lab/internal/config"
	llmopenai "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadClient()
	if err != nil {
		log.Fatal(err)
	}

	mcpClient, err := mcpclient.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mcpClient.Close()

	llmClient := llmopenai.New(cfg.AgentModel)

	agentClient := agent.New(mcpClient, llmClient)

	response, err := agentClient.RunPrompt(ctx, "task-review", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Resposta:")
	fmt.Println(response)
}

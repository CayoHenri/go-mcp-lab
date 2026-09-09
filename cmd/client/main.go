package main

import (
	"context"
	"log"
	"os"

	"github.com/CayoHenri/go-mcp-lab/internal/agent"
	"github.com/CayoHenri/go-mcp-lab/internal/cli"
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

	console := cli.NewConsole(os.Stdin, os.Stdout)

	approver := cli.NewApprover(console)

	agentClient := agent.New(mcpClient, llmClient, approver)

	app := cli.New(mcpClient, agentClient, console)

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

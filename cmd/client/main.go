package main

import (
	"context"
	"fmt"
	"log"

	"github.com/CayoHenri/go-mcp-lab/internal/agent"
	llmopenai "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
)

func main() {
	ctx := context.Background()

	mcpClient, err := mcpclient.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mcpClient.Close()

	llmClient := llmopenai.New("gpt-5.6-luna")

	agentClient := agent.New(mcpClient, llmClient)

	question := "Quanto é 35 vezes 12?"

	response, err := agentClient.Run(ctx,question)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Resposta:")
	fmt.Println(response)
}

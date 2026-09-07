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

	_, err := config.LoadClient()
	if err != nil {
		log.Fatal(err)
	}

	mcpClient, err := mcpclient.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mcpClient.Close()

	llmClient := llmopenai.New("gpt-5.6-luna")

	agentClient := agent.New(mcpClient, llmClient)
	question := `
	Liste todas as tarefas existentes.
	Completo a ultima tarefa da lista.
	Em seguida, liste novamente todas as tarefas existentes.
	`

	response, err := agentClient.Run(ctx, question)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Resposta:")
	fmt.Println(response)
}

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

	question := `
	Crie uma tarefa chamada "Aprender MCP".
	Crie outra chamada "Aprender Agent Loop".

	Liste as tarefas.

	Depois conclua a tarefa "Aprender MCP".

	Por fim, liste novamente as tarefas e me diga
	qual está concluída e qual continua pendente.
`

	response, err := agentClient.Run(ctx, question)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Resposta:")
	fmt.Println(response)
}

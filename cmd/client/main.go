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

	resources, err := mcpClient.ListResources(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Resources disponíveis:")

	for _, resource := range resources {
		fmt.Printf(
			"- %s: %s\n",
			resource.URI,
			resource.Description,
		)
	}

	result, err :=
		mcpClient.ReadResource(
			ctx,
			"tasks://summary",
		)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Conteúdo:")

	for _, content := range result.Contents {
		fmt.Println(content.Text)
	}

	templates, err := mcpClient.ListResourceTemplates(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Resource Templates disponíveis:")

	for _, template := range templates {
		fmt.Printf("- %s: %s\n", template.URITemplate, template.Description)
	}

	res, err := mcpClient.ReadResource(ctx, "tasks://1")
	if err != nil {
		log.Fatal(err)
	}

	for _, content := range res.Contents {
		fmt.Println(content.Text)
	}

	llmClient := llmopenai.New("gpt-5.6-luna")

	agentClient := agent.New(mcpClient, llmClient)
	question := `
	Liste todas as tarefas existentes.
	`

	response, err := agentClient.Run(ctx, question)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Resposta:")
	fmt.Println(response)
}

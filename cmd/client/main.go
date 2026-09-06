package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "go-mcp-lab-client",
			Version: "v0.1.0",
		},
		nil,
	)

	session, err := client.Connect(
		ctx,
		&mcp.CommandTransport{
			Command: exec.Command("go", "run", "./cmd/server"),
		},
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	fmt.Println("Cliente conectado ao MCP Server")

	result, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tools disponíveis:")

	for _, tool := range result.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	fmt.Println()

	callResult, err := session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "greet",
			Arguments: map[string]any{
				"name": "Caio",
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Resultado da Tool:")

	for _, content := range callResult.Content {
		fmt.Println(content)
	}
}
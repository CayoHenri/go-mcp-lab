package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/CayoHenri/go-mcp-lab/internal/mcp/tools"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name: "go-mcp-lab", Version: "v0.1.0",
	}, nil)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "greet",
			Description: "Cumprimenta uma pessoa pelo nome",
		},
		tools.Greet,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "calculate",
			Description: "Realiza operações matemáticas básicas entre dois números",
		},
		tools.Calculate,
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

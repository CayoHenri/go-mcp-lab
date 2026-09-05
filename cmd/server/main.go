package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/CayoHenri/go-mcp-lab/internal/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// 1. Configurando metadados avançados do servidor
	serverInfo := &mcp.Implementation{
		Name:        "go-mcp-lab",
		Version:     "v0.1.0",
		Title:       "Laboratório Go MCP",
		Description: "Ambiente de testes para ferramentas integradas a LLMs",
	}

	// 2. Definindo opções estruturadas do servidor
	serverOpts := &mcp.ServerOptions{
		Instructions: "Sempre liste os recursos disponíveis antes de executar qualquer ferramenta de escrita.",
		Logger:       slog.New(slog.NewTextHandler(os.Stdout, nil)),
		PageSize:     50,
		InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
			// Lógica executada logo após o handshake inicial com o LLM Client
			slog.Info("Cliente MCP conectado e inicializado com sucesso!")
		},
	}

	// 3. Inicializando o servidor
	server := mcp.NewServer(serverInfo, serverOpts)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "greet",
			Description: "Cumprimenta uma pessoa pelo nome",
		},
		tools.Greet,
	)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

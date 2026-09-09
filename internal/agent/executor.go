package agent

import (
	"context"
	"encoding/json"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (a *Agent) executeFunction(ctx context.Context, call llm.FunctionCall) (string, error) {
	fmt.Printf("Função selecionada: %s\n", call.Name)

	fmt.Printf("Argumentos: %+v\n", call.Arguments)

	if call.Name == "read_resource" {
		result, err := a.readResource(ctx, call.Arguments)
		if err != nil {
			return "", err
		}

		fmt.Printf("Resultado MCP Resource: %s\n", result)

		return result, nil
	}

	result, err := a.mcp.CallTool(ctx, call.Name, call.Arguments)
	if err != nil {
		return "", fmt.Errorf("executando MCP tool: %w", err)
	}

	resultText, err := serializeToolResult(result)
	if err != nil {
		return "", err
	}

	fmt.Printf("Resultado MCP Tool: %s\n", resultText)

	return resultText, nil
}

// readResource lê um Resource do MCP Server e retorna seu conteúdo serializado em JSON.
func (a *Agent) readResource(ctx context.Context, arguments map[string]any) (string, error) {
	uriValue, ok := arguments["uri"]

	if !ok {
		return "", fmt.Errorf("argumento uri não informado")
	}

	uri, ok := uriValue.(string)

	if !ok || uri == "" {
		return "", fmt.Errorf("argumento uri inválido")
	}

	result, err := a.mcp.ReadResource(ctx, uri)
	if err != nil {
		return "", err
	}

	return serializeResourceResult(result)
}

// serializeToolResult serializa o resultado de uma chamada de Tool em uma string JSON.
func serializeToolResult(result *mcp.CallToolResult) (string, error) {
	if result.IsError {
		data, err := json.Marshal(
			map[string]any{
				"error":   true,
				"content": result.Content,
			},
		)
		if err != nil {
			return "", err
		}

		return string(data), nil
	}

	if result.StructuredContent != nil {
		data, err := json.Marshal(result.StructuredContent)
		if err != nil {
			return "", err
		}

		return string(data), nil
	}

	data, err := json.Marshal(result.Content)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// serializeResourceResult serializa o resultado de uma leitura de Resource em uma string JSON.
func serializeResourceResult(result *mcp.ReadResourceResult) (string, error) {
	if result == nil {
		return "", fmt.Errorf("resultado de resource vazio")
	}

	data, err := json.Marshal(result.Contents)
	if err != nil {
		return "", fmt.Errorf("serializando resource: %w", err)
	}

	return string(data), nil
}

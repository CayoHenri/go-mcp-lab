package agent

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

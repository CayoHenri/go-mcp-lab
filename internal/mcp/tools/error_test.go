package tools

import (
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolError(t *testing.T) {
	err := errors.New("tarefa não encontrada")

	result := toolError(err)
	if result == nil {
		t.Fatal("esperado CallToolResult")
	}

	if !result.IsError {
		t.Error("esperado IsError true")
	}

	if len(result.Content) != 1 {
		t.Fatalf("esperado 1 conteúdo, recebido %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("esperado *mcp.TextContent, recebido %T", result.Content[0])
	}

	if textContent.Text != "tarefa não encontrada" {
		t.Errorf("mensagem inesperada: %q", textContent.Text)
	}
}

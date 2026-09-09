package agent

import (
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// promptToText converte um prompt do MCP em uma string de texto legível.
func promptToText(prompt *mcp.GetPromptResult) (string, error) {
	if prompt == nil {
		return "", fmt.Errorf("prompt MCP vazio")
	}

	var builder strings.Builder

	if prompt.Description != "" {
		builder.WriteString(prompt.Description)
		builder.WriteString("\n\n")
	}

	for _, message := range prompt.Messages {
		textContent, ok := message.Content.(*mcp.TextContent)

		if !ok {
			continue
		}

		if textContent.Text == "" {
			continue
		}

		fmt.Fprintf(&builder, "%s:\n%s\n\n", message.Role, textContent.Text)
	}

	result := strings.TrimSpace(builder.String())

	if result == "" {
		return "", fmt.Errorf("prompt MCP não possui conteúdo textual")
	}

	return result, nil
}

package openai

import (
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	sdk "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// resourceReaderTool creates a ToolUnionParam for reading MCP resources.
// It generates a description that includes the available resources and templates.
func resourceReaderTool(resources []*mcp.Resource, templates []*mcp.ResourceTemplate) responses.ToolUnionParam {
	description := buildResourceDescription(resources, templates)

	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        "read_resource",
			Description: sdk.String(description),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"uri": map[string]any{
						"type":        "string",
						"description": "URI MCP do resource que deve ser lido",
					},
				},
				"required": []string{
					"uri",
				},
				"additionalProperties": false,
			},
		},
	}
}

// buildResourceDescription builds a description string for the given resources and templates.
// It includes a general description of the resource handler, followed by a list of available resources and templates.
// The resulting string is formatted for readability.
func buildResourceDescription(resources []*mcp.Resource, templates []*mcp.ResourceTemplate) string {
	var builder strings.Builder

	builder.WriteString("Lê um MCP Resource. Use quando precisar consultar informações disponíveis no servidor.")

	if len(resources) > 0 {
		builder.WriteString("\n\nResources disponíveis:")
		for _, resource := range resources {
			fmt.Fprintf(&builder, "\n- %s: %s", resource.URI, resource.Description)
		}
	}

	if len(templates) > 0 {
		builder.WriteString("\n\nResource templates disponíveis:")
		for _, template := range templates {
			fmt.Fprintf(&builder, "\n- %s: %s", template.URITemplate, template.Description)
		}
	}

	return builder.String()
}

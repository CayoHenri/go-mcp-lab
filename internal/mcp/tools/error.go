package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

func toolError(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: err.Error(),
			},
		},
	}
}

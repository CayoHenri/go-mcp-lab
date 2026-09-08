/*
A responsabilidade deste pacote é:
MCP Client
 ├── conectar
 ├── listar Tools
 ├── chamar Tool
 └── fechar sessão
*/

package client

import (
	"context"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client struct {
	session *mcp.ClientSession
}

func New(ctx context.Context) (*Client, error) {
	mcpClient := mcp.NewClient(
		&mcp.Implementation{
			Name:    "go-mcp-lab-client",
			Version: "v0.1.0",
		},
		nil,
	)

	session, err := mcpClient.Connect(
		ctx,
		&mcp.CommandTransport{
			Command: exec.Command(
				"go",
				"run",
				"./cmd/server",
			),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Client{session: session}, nil
}

func (c *Client) Close() error {
	return c.session.Close()
}

// ListTools lista as Tools disponíveis no MCP Server.
func (c *Client) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	return c.session.ListTools(ctx, nil)
}

// CallTool chama uma Tool no MCP Server.
func (c *Client) CallTool(ctx context.Context, name string, arguments map[string]any) (*mcp.CallToolResult, error) {
	return c.session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name:      name,
			Arguments: arguments,
		},
	)
}

// ListResources lista os Resources disponíveis no MCP Server.
func (c *Client) ListResources(ctx context.Context) ([]*mcp.Resource, error) {
	resources := make([]*mcp.Resource, 0)

	for resource, err := range c.session.Resources(ctx, nil) {
		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

// ReadResource lê um Resource do MCP Server.
func (c *Client) ReadResource(ctx context.Context, uri string) (*mcp.ReadResourceResult, error) {
	return c.session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: uri,
	})
}

// ListResourceTemplates lista os ResourceTemplates disponíveis no MCP Server.
func (c *Client) ListResourceTemplates(ctx context.Context) ([]*mcp.ResourceTemplate, error) {
	templates := make([]*mcp.ResourceTemplate, 0)

	for template, err := range c.session.ResourceTemplates(ctx, nil) {
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	return templates, nil
}

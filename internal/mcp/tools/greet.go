package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetInput struct {
	Name string `json:"name" jsonschema:"nome da pessoa que será cumprimentada"`
}

type GreetOutput struct {
	Greeting string `json:"greeting" jsonschema:"mensagem de saudação"`
}

func Greet(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GreetInput,
) (*mcp.CallToolResult, GreetOutput, error) {
	return nil, GreetOutput{
		Greeting: fmt.Sprintf("Olá, %s!", input.Name),
	}, nil
}

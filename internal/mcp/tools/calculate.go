package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CalculateInput struct {
	Operation string  `json:"operation" jsonschema:"operação matemática: add, subtract, multiply ou divide"`
	A         float64 `json:"a" jsonschema:"primeiro número"`
	B         float64 `json:"b" jsonschema:"segundo número"`
}

type CalculateOutput struct {
	Result float64 `json:"result"`
}

func Calculate(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CalculateInput,
) (*mcp.CallToolResult, CalculateOutput, error) {
	var result float64

	switch input.Operation {
	case "add":
		result = input.A + input.B

	case "subtract":
		result = input.A - input.B

	case "multiply":
		result = input.A * input.B

	case "divide":
		if input.B == 0 {
			return nil, CalculateOutput{}, fmt.Errorf("não é possível dividir por zero")
		}

		result = input.A / input.B

	default:
		return nil, CalculateOutput{}, fmt.Errorf("operação inválida: %s", input.Operation)
	}

	return nil, CalculateOutput{Result: result}, nil
}

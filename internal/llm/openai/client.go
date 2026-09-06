/*
OpenAI Client
     │
     ├── converte MCP Tools
     │       ↓
     │   LLM Tools
     │
     ├── envia pergunta
     │
     ├── identifica Tool Call
     │
     └── devolve resultado ao modelo
*/

package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	sdk "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]any
}

type Client struct {
	client sdk.Client
	model  string
}

func New(model string) *Client {
	return &Client{
		client: sdk.NewClient(),
		model:  model,
	}
}

// Ask envia uma pergunta para o LLM e retorna a resposta.
func (c *Client) Ask(ctx context.Context, question string, tools []*mcp.Tool) (*responses.Response, error) {
	llmTools := make([]responses.ToolUnionParam, 0, len(tools))

	for _, tool := range tools {
		parameters, err := schemaToMap(tool.InputSchema)
		if err != nil {
			return nil, fmt.Errorf("convertendo schema da tool %s: %w", tool.Name, err)
		}

		llmTools = append(
			llmTools,
			responses.ToolUnionParam{
				OfFunction: &responses.FunctionToolParam{
					Name:        tool.Name,
					Description: sdk.String(tool.Description),
					Parameters:  parameters,
				},
			},
		)
	}

	response, err := c.client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: c.model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: sdk.String(question),
			},
			Tools: llmTools,
		},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ExtractToolCall verifica se a resposta do LLM contém uma chamada de Tool e retorna os detalhes da chamada.
func (c *Client) ExtractToolCall(response *responses.Response) (*ToolCall, bool, error) {
	for _, item := range response.Output {
		if item.Type != "function_call" {
			continue
		}

		call := item.AsFunctionCall()

		var arguments map[string]any

		if err := json.Unmarshal([]byte(call.Arguments), &arguments); err != nil {
			return nil, false, err
		}

		return &ToolCall{
			ID:        call.CallID,
			Name:      call.Name,
			Arguments: arguments,
		}, true, nil
	}

	return nil, false, nil
}

// SendToolResult envia o resultado da execução de uma Tool de volta para o LLM.
func (c *Client) SendToolResult(ctx context.Context, previousResponseID, callID, result string) (*responses.Response, error) {
	return c.client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model:              c.model,
			PreviousResponseID: sdk.String(previousResponseID),
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: []responses.ResponseInputItemUnionParam{
					{
						OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
							CallID: sdk.String(callID),
							Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
								OfString: sdk.String(result),
							},
						},
					},
				},
			},
		},
	)
}

// schemaToMap converte um schema de Tool em um mapa de strings para qualquer tipo.
func schemaToMap(schema any) (map[string]any, error) {
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	var result map[string]any

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

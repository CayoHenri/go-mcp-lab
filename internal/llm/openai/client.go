/*
OpenAI Client
 │
 ├── Ask
 │
 ├── MCP Tool → LLM Tool
 │
 ├── extrai todas Tool Calls
 │
 └── devolve todos Tool Results
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

type FunctionCall struct {
	ID        string
	Name      string
	Arguments map[string]any
}

type ToolResult struct {
	CallID string
	Output string
}

type Client struct {
	client sdk.Client
	model  string
}

type Resource struct {
	URI         string
	Name        string
	Description string
}

type ResourceTemplate struct {
	URITemplate string
	Name        string
	Description string
}

func New(model string) *Client {
	return &Client{
		client: sdk.NewClient(),
		model:  model,
	}
}

// Ask envia uma pergunta para o LLM e retorna a resposta.
func (c *Client) Ask(
	ctx context.Context,
	question string,
	tools []*mcp.Tool,
	resources []*mcp.Resource,
	templates []*mcp.ResourceTemplate,
) (*responses.Response, error) {
	llmTools, err := convertTools(tools, resources, templates)
	if err != nil {
		return nil, err
	}

	return c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: sdk.String(question),
		},
		Tools: llmTools,
	})
}

// // ExtractFunctionCalls extrai todas as Tool Calls da resposta do LLM.
func (c *Client) ExtractFunctionCalls(response *responses.Response) ([]FunctionCall, error) {
	calls := make([]FunctionCall, 0)

	for _, item := range response.Output {
		if item.Type != "function_call" {
			continue
		}

		call := item.AsFunctionCall()

		var arguments map[string]any

		if err := json.Unmarshal([]byte(call.Arguments), &arguments); err != nil {
			return nil, fmt.Errorf("decodificando argumentos da tool %s: %w", call.Name, err)
		}

		calls = append(
			calls,
			FunctionCall{
				ID:        call.CallID,
				Name:      call.Name,
				Arguments: arguments,
			},
		)
	}

	return calls, nil
}

// SendToolResults envia os resultados das chamadas de ferramentas (Tool Results) de volta para o LLM, permitindo que ele continue a interação com base nos resultados obtidos.
func (c *Client) SendToolResults(
	ctx context.Context,
	previousResponseID string,
	results []ToolResult,
	tools []*mcp.Tool,
	resources []*mcp.Resource,
	templates []*mcp.ResourceTemplate,
) (*responses.Response, error) {
	llmTools, err := convertTools(tools, resources, templates)
	if err != nil {
		return nil, err
	}

	inputs := make([]responses.ResponseInputItemUnionParam, 0, len(results))
	for _, result := range results {
		inputs = append(
			inputs,
			responses.ResponseInputItemUnionParam{
				OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
					CallID: sdk.String(result.CallID),
					Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
						OfString: sdk.String(result.Output),
					},
				},
			},
		)
	}

	return c.client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model:              c.model,
			PreviousResponseID: sdk.String(previousResponseID),
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: inputs,
			},
			Tools: llmTools, // Inclui as ferramentas (Tools) para que o LLM possa continuar a interação com elas.
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

func convertTools(
	tools []*mcp.Tool,
	resources []*mcp.Resource,
	templates []*mcp.ResourceTemplate,
) ([]responses.ToolUnionParam, error) {
	result := make([]responses.ToolUnionParam, 0, len(tools)+1)

	for _, tool := range tools {
		parameters, err := schemaToMap(tool.InputSchema)
		if err != nil {
			return nil, fmt.Errorf("convertendo schema da tool %s: %w", tool.Name, err)
		}

		result = append(result, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        tool.Name,
				Description: sdk.String(tool.Description),
				Parameters:  parameters,
			},
		})
	}

	if len(resources) > 0 || len(templates) > 0 {
		result = append(result, resourceReaderTool(
			resources,
			templates,
		))
	}

	return result, nil
}

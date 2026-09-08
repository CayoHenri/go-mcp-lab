package agent

import (
	"context"
	"encoding/json"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxIterations = 10

type Agent struct {
	mcp *mcpclient.Client
	llm *llm.Client
}

func New(mcp *mcpclient.Client, llmClient *llm.Client) *Agent {
	return &Agent{
		mcp: mcp,
		llm: llmClient,
	}
}

// Run executa o fluxo de interação entre o LLM e o MCP Server.
func (a *Agent) Run(ctx context.Context, question string) (string, error) {
	toolsResult, err := a.mcp.ListTools(ctx)
	if err != nil {
		return "", fmt.Errorf("listando MCP tools: %w", err)
	}

	resources, err := a.mcp.ListResources(ctx)
	if err != nil {
		return "", fmt.Errorf("listando MCP resources: %w", err)
	}

	templates, err := a.mcp.ListResourceTemplates(ctx)
	if err != nil {
		return "", fmt.Errorf("listando MCP resource templates: %w", err)
	}

	response, err := a.llm.Ask(ctx, question, toolsResult.Tools, resources, templates)
	if err != nil {
		return "", fmt.Errorf("consultando LLM: %w", err)
	}

	for iteration := 1; iteration <= maxIterations; iteration++ {
		toolCalls, err := a.llm.ExtractFunctionCalls(response)
		if err != nil {
			return "", fmt.Errorf("interpretando tool calls: %w", err)
		}

		if len(toolCalls) == 0 {
			return response.OutputText(), nil
		}

		fmt.Printf("\n[Iteração %d]\n", iteration)

		toolResults := make([]llm.ToolResult, 0, len(toolCalls))

		for _, toolCall := range toolCalls {
			fmt.Printf("Tool selecionada: %s\n", toolCall.Name)

			fmt.Printf("Argumentos: %+v\n", toolCall.Arguments)

			if toolCall.Name == "read_resource" {
				resultText, err := a.readResource(ctx, toolCall.Arguments)
				if err != nil {
					return "", fmt.Errorf("lendo MCP resource: %w", err)
				}

				toolResults = append(toolResults, llm.ToolResult{
					CallID: toolCall.ID,
					Output: resultText,
				})

				continue
			}

			result, err := a.mcp.CallTool(ctx, toolCall.Name, toolCall.Arguments)
			if err != nil {
				return "", fmt.Errorf("executando MCP tool %s: %w", toolCall.Name, err)
			}

			resultText, err := serializeToolResult(result)
			if err != nil {
				return "", fmt.Errorf("serializando resultado da tool %s: %w", toolCall.Name, err)
			}

			fmt.Printf("Resultado MCP: %s\n", resultText)

			toolResults = append(
				toolResults,
				llm.ToolResult{
					CallID: toolCall.ID,
					Output: resultText,
				},
			)
		}

		response, err = a.llm.SendToolResults(ctx, response.ID, toolResults, toolsResult.Tools, resources, templates)
		if err != nil {
			return "", fmt.Errorf("enviando resultados ao LLM: %w", err)
		}
	}

	return "", fmt.Errorf("limite máximo de %d iterações atingido", maxIterations)
}

func (a *Agent) readResource(ctx context.Context, arguments map[string]any) (string, error) {
	uriValue, ok := arguments["uri"]

	if !ok {
		return "", fmt.Errorf("argumento uri não informado")
	}

	uri, ok := uriValue.(string)

	if !ok || uri == "" {
		return "", fmt.Errorf("argumento uri inválido")
	}

	result, err := a.mcp.ReadResource(ctx, uri)
	if err != nil {
		return "", err
	}

	return serializeResourceResult(result)
}

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

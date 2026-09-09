package agent

import (
	"context"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
)

const maxIterations = 10

func (a *Agent) runLoop(ctx context.Context, input string) (string, error) {
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

	response, err := a.llm.Ask(ctx, input, toolsResult.Tools, resources, templates, a.session.PreviousResponseID())
	if err != nil {
		return "", fmt.Errorf("consultando LLM: %w", err)
	}

	for iteration := 1; iteration <= maxIterations; iteration++ {
		functionCalls, err := a.llm.ExtractFunctionCalls(response)
		if err != nil {
			return "", fmt.Errorf("interpretando function calls: %w", err)
		}

		if len(functionCalls) == 0 {
			a.session.SetPreviousResponseID(response.ID)
			return response.OutputText(), nil
		}

		fmt.Printf("\n[Iteração %d]\n", iteration)

		results := make([]llm.ToolResult, 0, len(functionCalls))

		for _, call := range functionCalls {
			resultText, err := a.executeFunction(ctx, call)
			if err != nil {
				return "", fmt.Errorf("executando função %s: %w", call.Name, err)
			}

			results = append(results, llm.ToolResult{
				CallID: call.ID,
				Output: resultText,
			})
		}

		response, err = a.llm.SendToolResults(ctx, response.ID, results, toolsResult.Tools, resources, templates)
		if err != nil {
			return "", fmt.Errorf("enviando resultados ao LLM: %w", err)
		}
	}

	return "", fmt.Errorf("limite máximo de %d iterações atingido", maxIterations)
}

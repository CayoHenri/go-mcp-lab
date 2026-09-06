package agent

import (
	"context"
	"encoding/json"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
)

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

// Run executa o fluxo completo do agente: envia a pergunta para o LLM, identifica se há uma Tool Call, executa a Tool no MCP Server e envia o resultado de volta para o LLM.
func (a *Agent) Run(ctx context.Context, question string) (string, error) {
	toolsResult, err := a.mcp.ListTools(ctx)
	if err != nil {
		return "", fmt.Errorf("listando MCP tools: %w", err)
	}

	response, err := a.llm.Ask(ctx, question, toolsResult.Tools)
	if err != nil {
		return "", fmt.Errorf("consultando LLM: %w", err)
	}

	toolCall, hasToolCall, err := a.llm.ExtractToolCall(response)
	if err != nil {
		return "", fmt.Errorf("interpretando tool call: %w", err)
	}

	if !hasToolCall {
		return response.OutputText(), nil
	}

	fmt.Println("LLM selecionou:")
	fmt.Printf("Tool: %s\n", toolCall.Name)
	fmt.Printf("Argumentos: %+v\n", toolCall.Arguments)

	result, err := a.mcp.CallTool(ctx, toolCall.Name, toolCall.Arguments)
	if err != nil {
		return "", fmt.Errorf("executando MCP tool %s: %w", toolCall.Name, err)
	}

	resultJSON, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return "", fmt.Errorf("serializando resultado da tool: %w", err)
	}

	fmt.Printf("Resultado MCP: %s\n", string(resultJSON))

	finalResponse, err := a.llm.SendToolResult(ctx, response.ID, toolCall.ID, string(resultJSON))
	if err != nil {
		return "", fmt.Errorf("enviando resultado ao LLM: %w", err)
	}

	return finalResponse.OutputText(), nil
}

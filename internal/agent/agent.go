package agent

import (
	"context"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
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

// Run executa o fluxo de interação entre o LLM e o MCP Server, processando a pergunta fornecida e retornando a resposta final.
func (a *Agent) Run(ctx context.Context, question string) (string, error) {
	return a.runLoop(ctx, question)
}

// RunPrompt executa um prompt específico do MCP, processando os argumentos fornecidos e retornando a resposta final.
func (a *Agent) RunPrompt(ctx context.Context, promptName string, arguments map[string]string) (string, error) {
	prompt, err := a.mcp.GetPrompt(ctx, promptName, arguments)
	if err != nil {
		return "", fmt.Errorf("obtendo MCP prompt %s: %w", promptName, err)
	}

	input, err := promptToText(prompt)
	if err != nil {
		return "", fmt.Errorf("convertendo MCP prompt %s: %w", promptName, err)
	}

	return a.runLoop(ctx, input)
}

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

	response, err := a.llm.Ask(ctx, input, toolsResult.Tools, resources, templates)
	if err != nil {
		return "", fmt.Errorf("consultando LLM: %w", err)
	}

	for iteration := 1; iteration <= maxIterations; iteration++ {
		functionCalls, err := a.llm.ExtractFunctionCalls(response)
		if err != nil {
			return "", fmt.Errorf("interpretando function calls: %w", err)
		}

		if len(functionCalls) == 0 {
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

func (a *Agent) executeFunction(ctx context.Context, call llm.FunctionCall) (string, error) {
	fmt.Printf("Função selecionada: %s\n", call.Name)

	fmt.Printf("Argumentos: %+v\n", call.Arguments)

	if call.Name == "read_resource" {
		result, err := a.readResource(ctx, call.Arguments)
		if err != nil {
			return "", err
		}

		fmt.Printf("Resultado MCP Resource: %s\n", result)

		return result, nil
	}

	result, err := a.mcp.CallTool(ctx, call.Name, call.Arguments)
	if err != nil {
		return "", fmt.Errorf("executando MCP tool: %w", err)
	}

	resultText, err := serializeToolResult(result)
	if err != nil {
		return "", err
	}

	fmt.Printf("Resultado MCP Tool: %s\n", resultText)

	return resultText, nil
}

// readResource lê um Resource do MCP Server e retorna seu conteúdo serializado em JSON.
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

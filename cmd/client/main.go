package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	ctx := context.Background()

	// 1. Cria o MCP Client.
	mcpClient := mcp.NewClient(
		&mcp.Implementation{
			Name:    "go-mcp-lab-client",
			Version: "v0.1.0",
		},
		nil,
	)

	// 2. Inicia e conecta ao nosso MCP Server.
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
		log.Fatal(err)
	}
	defer session.Close()

	fmt.Println("Cliente conectado ao MCP Server")

	// 3. Descobre as Tools disponíveis.
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tools descobertas:")

	for _, tool := range toolsResult.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// 4. Converte as Tools MCP para Tools entendidas pelo LLM.
	llmTools := make([]responses.ToolUnionParam, 0, len(toolsResult.Tools))

	for _, tool := range toolsResult.Tools {
		parameters, err := schemaToMap(tool.InputSchema)
		if err != nil {
			log.Fatalf("erro convertendo schema da tool %s: %v", tool.Name, err)
		}

		llmTools = append(
			llmTools,
			responses.ToolUnionParam{
				OfFunction: &responses.FunctionToolParam{
					Name:        tool.Name,
					Description: openai.String(tool.Description),
					Parameters:  parameters,
				},
			},
		)
	}

	// 5. Cria o cliente da OpenAI.
	openAIClient := openai.NewClient()

	question := "Cumprimente o Carlos"

	fmt.Println()
	fmt.Println("Usuário:", question)

	// 6. Envia a pergunta + Tools para o modelo.
	response, err := openAIClient.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: "gpt-5.6-luna",

			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(question),
			},

			Tools: llmTools,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// 7. Verifica se o modelo decidiu chamar alguma Tool.
	for _, item := range response.Output {
		if item.Type != "function_call" {
			continue
		}

		toolCall := item.AsFunctionCall()

		fmt.Println()
		fmt.Println("LLM decidiu chamar:")
		fmt.Println("Tool:", toolCall.Name)
		fmt.Println("Argumentos:", toolCall.Arguments)

		var arguments map[string]any

		if err := json.Unmarshal([]byte(toolCall.Arguments),&arguments); err != nil {
			log.Fatal(err)
		}

		// 8. Executa a Tool de verdade através do MCP.
		toolResult, err := session.CallTool(
			ctx,
			&mcp.CallToolParams{
				Name:      toolCall.Name,
				Arguments: arguments,
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		// 9. Serializa o resultado da Tool.
		resultJSON, err := json.Marshal(toolResult.StructuredContent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Resultado MCP:", string(resultJSON))

		// 10. Devolve o resultado da Tool ao LLM.
		finalResponse, err := openAIClient.Responses.New(
			ctx,
			responses.ResponseNewParams{
				Model:              "gpt-5.6-luna",
				PreviousResponseID: openai.String(response.ID),

				Input: responses.ResponseNewParamsInputUnion{
					OfInputItemList: []responses.ResponseInputItemUnionParam{
						{
							OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
								CallID: openai.String(toolCall.CallID),

								Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
									OfString: openai.String(string(resultJSON)),
								},
							},
						},
					},
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println()
		fmt.Println("Resposta final:")
		fmt.Println(finalResponse.OutputText())

		return
	}

	// Caso o modelo não precise de Tool.
	fmt.Println()
	fmt.Println("Resposta direta:")
	fmt.Println(response.OutputText())
}

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

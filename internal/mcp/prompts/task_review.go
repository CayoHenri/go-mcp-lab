package prompts

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const TaskReviewName = "task-review"

func TaskReview(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Orientações para revisar o estado atual das tarefas",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: `
Analise o estado atual das tarefas.

Consulte os resources MCP disponíveis quando precisar obter os dados.

Apresente:
- quantidade total de tarefas;
- quantidade de tarefas concluídas;
- quantidade de tarefas pendentes;
- quais tarefas ainda estão pendentes;
- um resumo objetivo do estado atual.

Não crie, conclua ou altere tarefas durante esta análise.
`,
				},
			},
		},
	}, nil
}

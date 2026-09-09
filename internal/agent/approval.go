package agent

import (
	"context"
	"encoding/json"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
)

type Approver interface {
	Approve(ctx context.Context, call llm.FunctionCall) (bool, error)
}

func requiresApproval(call llm.FunctionCall) bool {
	switch call.Name {
	case "create_task", "complete_task":
		return true

	default:
		return false
	}
}

func (a *Agent) requestApproval(ctx context.Context, call llm.FunctionCall) (bool, error) {
	if a.approver == nil {
		return false, nil
	}

	return a.approver.Approve(ctx, call)
}

func toolRejectedResult(call llm.FunctionCall) string {
	result := map[string]any{
		"success":  false,
		"rejected": true,
		"message":  fmt.Sprintf("Execução da função %s não autorizada pelo usuário.", call.Name),
	}

	data, err := json.Marshal(result)
	if err != nil {
		return `{"success":false,"rejected":true}`
	}

	return string(data)
}

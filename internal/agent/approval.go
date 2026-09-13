package agent

import (
	"context"
	"encoding/json"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
)

type ApprovalRequest struct {
	Call            llm.FunctionCall
	PermissionLevel PermissionLevel
}

type Approver interface {
	Approve(ctx context.Context, request ApprovalRequest) (bool, error)
}

func functionRejectedResult(call llm.FunctionCall) string {
	result := map[string]any{
		"success":  false,
		"rejected": true,
		"function": call.Name,
		"message":  fmt.Sprintf("Execução da função %s não autorizada.", call.Name),
	}

	data, err := json.Marshal(result)
	if err != nil {
		return `{"success":false,"rejected":true}`
	}

	return string(data)
}

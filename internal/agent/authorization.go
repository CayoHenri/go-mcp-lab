package agent

import (
	"context"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
)

func (a *Agent) authorize(ctx context.Context, call llm.FunctionCall) (bool, error) {
	level := a.policy.Level(call.Name)

	switch level {
	case PermissionRead:
		return true, nil

	case PermissionWrite, PermissionDestructive:
		if a.approver == nil {
			return false, nil
		}

		return a.approver.Approve(ctx, ApprovalRequest{
			Call:            call,
			PermissionLevel: level,
		})
	case PermissionUnknown:
		return false, fmt.Errorf("função %s não possui política de permissão definida", call.Name)

	default:
		return false, fmt.Errorf("nível de permissão inválido para função %s", call.Name)
	}
}

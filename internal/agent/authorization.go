package agent

import (
	"context"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (a *Agent) authorize(ctx context.Context, call llm.FunctionCall, tools []*mcp.Tool) (bool, error) {
	level := a.permissionLevel(call.Name, tools)

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
		return false, fmt.Errorf("não foi possível determinar a permissão da função %s", call.Name)

	default:
		return false, fmt.Errorf("nível de permissão inválido para %s", call.Name)
	}
}

func (a *Agent) permissionLevel(functionName string, tools []*mcp.Tool) PermissionLevel {
	if functionName == "read_resource" {
		return a.policy.InternalLevel(functionName)
	}

	tool := findTool(tools, functionName)
	if tool == nil {
		return PermissionUnknown
	}

	return a.policy.ToolLevel(tool)
}

func findTool(tools []*mcp.Tool, name string) *mcp.Tool {
	for _, tool := range tools {
		if tool.Name == name {
			return tool
		}
	}

	return nil
}

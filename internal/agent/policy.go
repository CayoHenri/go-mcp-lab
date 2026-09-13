package agent

import "github.com/modelcontextprotocol/go-sdk/mcp"

type PermissionPolicy struct{}

func NewPermissionPolicy() *PermissionPolicy {
	return &PermissionPolicy{}
}

func (p *PermissionPolicy) ToolLevel(tool *mcp.Tool) PermissionLevel {
	if tool == nil {
		return PermissionUnknown
	}

	annotations := tool.Annotations

	if annotations == nil {
		return PermissionUnknown
	}

	if annotations.ReadOnlyHint {
		return PermissionRead
	}

	if destructiveHint(annotations) {
		return PermissionDestructive
	}

	return PermissionWrite
}

func (p *PermissionPolicy) InternalLevel(functionName string) PermissionLevel {
	switch functionName {
	case "read_resource":
		return PermissionRead

	default:
		return PermissionUnknown
	}
}

// No protocolo MCP: destructiveHint default = true
// Ou seja, ausência não significa: false e sim true
// O SDK representa isso com: *bool justamente para diferenciar: nil || false || true
func destructiveHint(annotations *mcp.ToolAnnotations) bool {
	if annotations.DestructiveHint == nil {
		return true
	}

	return *annotations.DestructiveHint
}

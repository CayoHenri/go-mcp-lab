package agent

type PermissionPolicy struct {
	functions map[string]PermissionLevel
}

func NewPermissionPolicy() *PermissionPolicy {
	return &PermissionPolicy{
		functions: map[string]PermissionLevel{
			"greet":         PermissionRead,
			"calculate":     PermissionRead,
			"list_tasks":    PermissionRead,
			"read_resource": PermissionRead,
			"create_task":   PermissionWrite,
			"complete_task": PermissionWrite,
			"delete_task":   PermissionDestructive,
		},
	}
}

func (p *PermissionPolicy) Level(functionName string) PermissionLevel {
	level, ok := p.functions[functionName]
	if !ok {
		return PermissionUnknown
	}

	return level
}

func (p *PermissionPolicy) RequiresApproval(functionName string) bool {
	switch p.Level(functionName) {
	case PermissionWrite, PermissionDestructive:
		return true

	default:
		return false
	}
}

func (p *PermissionPolicy) IsAllowedWithoutApproval(functionName string) bool {
	return p.Level(functionName) == PermissionRead
}

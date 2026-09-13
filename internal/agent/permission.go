package agent

type PermissionLevel string

const (
	PermissionRead        PermissionLevel = "read"
	PermissionWrite       PermissionLevel = "write"
	PermissionDestructive PermissionLevel = "destructive"
	PermissionUnknown     PermissionLevel = "unknown"
)

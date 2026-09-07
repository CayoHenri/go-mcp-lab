package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/CayoHenri/go-mcp-lab/internal/task"
)

type CreateTaskInput struct {
	Title string `json:"title" jsonschema:"título da tarefa"`
}

type CreateTaskOutput struct {
	Task task.Task `json:"task"`
}

func CreateTask(store *task.Store) func(
	context.Context,
	*mcp.CallToolRequest,
	CreateTaskInput,
) (*mcp.CallToolResult, CreateTaskOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input CreateTaskInput,
	) (*mcp.CallToolResult, CreateTaskOutput, error) {
		created := store.Create(input.Title)

		return nil, CreateTaskOutput{Task: created}, nil
	}
}

type ListTasksInput struct{}

type ListTasksOutput struct {
	Tasks []task.Task `json:"tasks"`
}

func ListTasks(store *task.Store) func(
	context.Context,
	*mcp.CallToolRequest,
	ListTasksInput,
) (*mcp.CallToolResult, ListTasksOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input ListTasksInput,
	) (*mcp.CallToolResult, ListTasksOutput, error) {
		return nil, ListTasksOutput{
			Tasks: store.List(),
		}, nil
	}
}

type CompleteTaskInput struct {
	ID int `json:"id" jsonschema:"identificador da tarefa"`
}

type CompleteTaskOutput struct {
	Task task.Task `json:"task"`
}

func CompleteTask(store *task.Store) func(
	context.Context,
	*mcp.CallToolRequest,
	CompleteTaskInput,
) (*mcp.CallToolResult, CompleteTaskOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input CompleteTaskInput,
	) (*mcp.CallToolResult, CompleteTaskOutput, error) {
		completed, err := store.Complete(input.ID)
		if err != nil {
			return nil, CompleteTaskOutput{}, err
		}

		return nil, CompleteTaskOutput{Task: completed}, nil
	}
}

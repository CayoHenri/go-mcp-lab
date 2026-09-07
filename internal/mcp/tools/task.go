package tools

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	apptask "github.com/CayoHenri/go-mcp-lab/internal/application/task"
	domaintask "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

type CreateTaskInput struct {
	Title string `json:"title" jsonschema:"título da tarefa"`
}

type CreateTaskOutput struct {
	Task apptask.TaskOutput `json:"task"`
}

func CreateTask(
	service *apptask.Service,
) func(
	context.Context,
	*mcp.CallToolRequest,
	CreateTaskInput,
) (*mcp.CallToolResult, CreateTaskOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input CreateTaskInput,
	) (*mcp.CallToolResult, CreateTaskOutput, error) {
		created, err := service.Create(ctx, input.Title)
		if err != nil {
			return toolError(err), CreateTaskOutput{}, nil
		}

		return nil, CreateTaskOutput{
			Task: created,
		}, nil
	}
}

type ListTasksInput struct{}

type ListTasksOutput struct {
	Tasks []apptask.TaskOutput `json:"tasks"`
}

func ListTasks(
	service *apptask.Service,
) func(
	context.Context,
	*mcp.CallToolRequest,
	ListTasksInput,
) (*mcp.CallToolResult, ListTasksOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input ListTasksInput,
	) (*mcp.CallToolResult, ListTasksOutput, error) {
		tasks, err := service.List(ctx)
		if err != nil {
			return toolError(err), ListTasksOutput{}, nil
		}

		return nil, ListTasksOutput{
			Tasks: tasks,
		}, nil
	}
}

type CompleteTaskInput struct {
	ID int `json:"id" jsonschema:"identificador da tarefa"`
}

type CompleteTaskOutput struct {
	Task apptask.TaskOutput `json:"task"`
}

func CompleteTask(
	service *apptask.Service,
) func(
	context.Context,
	*mcp.CallToolRequest,
	CompleteTaskInput,
) (*mcp.CallToolResult, CompleteTaskOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input CompleteTaskInput,
	) (*mcp.CallToolResult, CompleteTaskOutput, error) {
		completed, err := service.Complete(ctx, input.ID)
		if err != nil {
			if errors.Is(err, domaintask.ErrNotFound) {
				return toolError(err), CompleteTaskOutput{}, nil
			}

			return toolError(err), CompleteTaskOutput{}, nil
		}

		return nil, CompleteTaskOutput{
			Task: completed,
		}, nil
	}
}

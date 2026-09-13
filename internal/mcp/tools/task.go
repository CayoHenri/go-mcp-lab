package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	apptask "github.com/CayoHenri/go-mcp-lab/internal/application/task"
	domaintask "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

// TaskService define o contrato de serviço para gerenciar tarefas
type TaskService interface {
	Create(ctx context.Context, title string) (apptask.TaskOutput, error)
	List(ctx context.Context) ([]apptask.TaskOutput, error)
	Complete(ctx context.Context, id int) (apptask.TaskOutput, error)
	Delete(ctx context.Context, id int) error
}

type CreateTaskInput struct {
	Title string `json:"title" jsonschema:"título da tarefa"`
}

type CreateTaskOutput struct {
	Task apptask.TaskOutput `json:"task"`
}

func CreateTask(service TaskService) func(
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

func ListTasks(service TaskService) func(
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

func CompleteTask(service TaskService) func(
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

type DeleteTaskInput struct {
	ID int `json:"id" jsonschema:"identificador da tarefa que será excluída"`
}

type DeleteTaskOutput struct {
	ID      int    `json:"id"`
	Deleted bool   `json:"deleted"`
	Message string `json:"message"`
}

func DeleteTask(service TaskService) func(
	context.Context,
	*mcp.CallToolRequest,
	DeleteTaskInput,
) (*mcp.CallToolResult, DeleteTaskOutput, error) {
	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input DeleteTaskInput,
	) (*mcp.CallToolResult, DeleteTaskOutput, error) {
		if err := service.Delete(ctx, input.ID); err != nil {
			return toolError(err), DeleteTaskOutput{}, nil
		}

		return nil, DeleteTaskOutput{
			ID:      input.ID,
			Deleted: true,
			Message: fmt.Sprintf("Tarefa %d excluída com sucesso", input.ID),
		}, nil
	}
}

package tools

import (
	"context"
	"errors"
	"testing"

	apptask "github.com/CayoHenri/go-mcp-lab/internal/application/task"
	domaintask "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

type taskServiceFake struct {
	createFn   func(context.Context, string) (apptask.TaskOutput, error)
	listFn     func(context.Context) ([]apptask.TaskOutput, error)
	completeFn func(context.Context, int) (apptask.TaskOutput, error)
}

func (s *taskServiceFake) Create(ctx context.Context, title string) (apptask.TaskOutput, error) {
	if s.createFn == nil {
		return apptask.TaskOutput{}, nil
	}

	return s.createFn(ctx, title)
}

func (s *taskServiceFake) List(ctx context.Context) ([]apptask.TaskOutput, error) {
	if s.listFn == nil {
		return nil, nil
	}

	return s.listFn(ctx)
}

func (s *taskServiceFake) Complete(ctx context.Context, id int) (apptask.TaskOutput, error) {
	if s.completeFn == nil {
		return apptask.TaskOutput{}, nil
	}

	return s.completeFn(ctx, id)
}

func TestCreateTask(t *testing.T) {
	var receivedTitle string

	service := &taskServiceFake{
		createFn: func(ctx context.Context, title string) (apptask.TaskOutput, error) {
			receivedTitle = title
			return apptask.TaskOutput{
				ID:        1,
				Title:     title,
				Completed: false,
			}, nil
		},
	}

	handler := CreateTask(service)

	result, output, err := handler(context.Background(), nil, CreateTaskInput{
		Title: "Estudar MCP",
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result != nil {
		t.Fatalf("esperado CallToolResult nil em caso de sucesso")
	}

	if receivedTitle != "Estudar MCP" {
		t.Errorf("esperado título %q, recebido %q", "Estudar MCP", receivedTitle)
	}

	if output.Task.ID != 1 {
		t.Errorf("esperado ID 1, recebido %d", output.Task.ID)
	}

	if output.Task.Title != "Estudar MCP" {
		t.Errorf("esperado título %q, recebido %q", "Estudar MCP", output.Task.Title)
	}

	if output.Task.Completed {
		t.Error("tarefa não deveria estar concluída")
	}
}

func TestCreateTaskReturnsToolError(t *testing.T) {
	expectedErr := errors.New("não foi possível criar a tarefa")

	service := &taskServiceFake{
		createFn: func(ctx context.Context, title string) (apptask.TaskOutput, error) {
			return apptask.TaskOutput{}, expectedErr
		},
	}

	handler := CreateTask(service)

	result, output, err := handler(context.Background(), nil, CreateTaskInput{
		Title: "Estudar MCP",
	})

	if err != nil {
		t.Fatalf("erro técnico não era esperado: %v", err)
	}

	if result == nil {
		t.Fatal("esperado CallToolResult")
	}

	if !result.IsError {
		t.Error("esperado IsError true")
	}

	if output != (CreateTaskOutput{}) {
		t.Errorf("esperado output vazio, recebido %+v", output)
	}

	if len(result.Content) == 0 {
		t.Fatal("esperado conteúdo de erro")
	}
}

func TestListTasks(t *testing.T) {
	service := &taskServiceFake{
		listFn: func(ctx context.Context) ([]apptask.TaskOutput, error) {
			return []apptask.TaskOutput{
				{
					ID:        1,
					Title:     "Estudar MCP",
					Completed: false,
				},
				{
					ID:        2,
					Title:     "Estudar agentes",
					Completed: true,
				},
			}, nil
		},
	}

	handler := ListTasks(service)

	result, output, err := handler(context.Background(), nil, ListTasksInput{})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result != nil {
		t.Fatal("esperado CallToolResult nil")
	}

	if len(output.Tasks) != 2 {
		t.Fatalf("esperadas 2 tarefas, recebidas %d", len(output.Tasks))
	}

	if output.Tasks[0].ID != 1 {
		t.Errorf("esperado primeiro ID 1, recebido %d", output.Tasks[0].ID)
	}

	if output.Tasks[1].ID != 2 {
		t.Errorf("esperado segundo ID 2, recebido %d", output.Tasks[1].ID)
	}

	if !output.Tasks[1].Completed {
		t.Error("segunda tarefa deveria estar concluída")
	}
}

func TestListTasksReturnsToolError(t *testing.T) {
	expectedErr := errors.New("erro ao listar tarefas")

	service := &taskServiceFake{
		listFn: func(ctx context.Context) ([]apptask.TaskOutput, error) {
			return nil, expectedErr
		},
	}

	handler := ListTasks(service)

	result, output, err := handler(context.Background(), nil, ListTasksInput{})

	if err != nil {
		t.Fatalf("erro técnico não era esperado: %v", err)
	}

	if result == nil {
		t.Fatal("esperado CallToolResult")
	}

	if !result.IsError {
		t.Error("esperado IsError true")
	}

	if output.Tasks != nil {
		t.Errorf("esperado Tasks nil, recebido %+v", output.Tasks)
	}
}

func TestCompleteTask(t *testing.T) {
	var receivedID int

	service := &taskServiceFake{
		completeFn: func(ctx context.Context, id int) (apptask.TaskOutput, error) {
			receivedID = id

			return apptask.TaskOutput{
				ID:        id,
				Title:     "Estudar MCP",
				Completed: true,
			}, nil
		},
	}

	handler := CompleteTask(service)

	result, output, err := handler(context.Background(), nil, CompleteTaskInput{
		ID: 1,
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if result != nil {
		t.Fatal("esperado CallToolResult nil")
	}

	if receivedID != 1 {
		t.Errorf("esperado ID 1 no service, recebido %d", receivedID)
	}

	if output.Task.ID != 1 {
		t.Errorf("esperado ID 1, recebido %d", output.Task.ID)
	}

	if !output.Task.Completed {
		t.Error("tarefa deveria estar concluída")
	}
}

func TestCompleteTaskNotFoundReturnsToolError(t *testing.T) {
	service := &taskServiceFake{
		completeFn: func(ctx context.Context, id int) (apptask.TaskOutput, error) {
			return apptask.TaskOutput{}, domaintask.ErrNotFound
		},
	}

	handler := CompleteTask(service)

	result, output, err := handler(context.Background(), nil, CompleteTaskInput{
		ID: 99,
	})

	if err != nil {
		t.Fatalf("erro técnico não era esperado: %v", err)
	}

	if result == nil {
		t.Fatal("esperado CallToolResult de erro")
	}

	if !result.IsError {
		t.Fatal("esperado IsError true")
	}

	if output != (CompleteTaskOutput{}) {
		t.Errorf("esperado output vazio, recebido %+v", output)
	}

	if len(result.Content) == 0 {
		t.Fatal("esperado conteúdo explicando o erro")
	}
}

func TestCompleteTaskInvalidIDReturnsToolError(t *testing.T) {
	service := &taskServiceFake{
		completeFn: func(ctx context.Context, id int) (apptask.TaskOutput, error) {
			return apptask.TaskOutput{}, domaintask.ErrInvalidID
		},
	}

	handler := CompleteTask(service)

	result, _, err := handler(context.Background(), nil, CompleteTaskInput{
		ID: 0,
	})

	if err != nil {
		t.Fatalf("erro técnico não era esperado: %v", err)
	}

	if result == nil {
		t.Fatal("esperado CallToolResult")
	}

	if !result.IsError {
		t.Error("esperado IsError true")
	}
}

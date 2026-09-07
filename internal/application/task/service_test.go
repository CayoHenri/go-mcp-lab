package task

import (
	"context"
	"errors"
	"testing"

	domain "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

type repositoryFake struct {
	createFn   func(context.Context, string) (domain.Task, error)
	listFn     func(context.Context) ([]domain.Task, error)
	completeFn func(context.Context, int) (domain.Task, error)
}

func (r *repositoryFake) Create(ctx context.Context, title string) (domain.Task, error) {
	if r.createFn == nil {
		return domain.Task{}, nil
	}

	return r.createFn(ctx, title)
}

func (r *repositoryFake) List(ctx context.Context) ([]domain.Task, error) {
	if r.listFn == nil {
		return nil, nil
	}

	return r.listFn(ctx)
}

func (r *repositoryFake) Complete(ctx context.Context, id int) (domain.Task, error) {
	if r.completeFn == nil {
		return domain.Task{}, nil
	}

	return r.completeFn(ctx, id)
}

func TestServiceCreate(t *testing.T) {
	repository := &repositoryFake{
		createFn: func(ctx context.Context, title string) (domain.Task, error) {
			task, err := domain.New(1, title)
			return task, err
		},
	}

	service := NewService(repository)

	output, err := service.Create(context.Background(), "Estudar MCP")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if output.ID != 1 {
		t.Errorf("esperado ID 1, recebido %d", output.ID)
	}

	if output.Title != "Estudar MCP" {
		t.Errorf("esperado título %q, recebido %q", "Estudar MCP", output.Title)
	}

	if output.Completed {
		t.Error("nova tarefa não deveria estar concluída")
	}
}

func TestServiceCreatePropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("erro ao criar tarefa")

	repository := &repositoryFake{
		createFn: func(ctx context.Context, title string) (domain.Task, error) {
			return domain.Task{}, expectedErr
		},
	}

	service := NewService(repository)

	output, err := service.Create(context.Background(), "Estudar MCP")

	if !errors.Is(err, expectedErr) {
		t.Fatalf("esperado erro %v, recebido %v", expectedErr, err)
	}

	if output != (TaskOutput{}) {
		t.Errorf("esperado TaskOutput vazio, recebido %+v", output)
	}
}

func TestServiceList(t *testing.T) {
	repository := &repositoryFake{
		listFn: func(ctx context.Context) ([]domain.Task, error) {
			task1, err := domain.New(1, "Estudar MCP")
			if err != nil {
				return nil, err
			}

			task2, err := domain.New(2, "Estudar agentes")
			if err != nil {
				return nil, err
			}
			task2.Complete()

			return []domain.Task{task1, task2}, nil
		},
	}

	service := NewService(repository)

	outputs, err := service.List(context.Background())

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(outputs) != 2 {
		t.Fatalf("esperadas 2 tarefas, recebidas %d", len(outputs))
	}

	if outputs[0].ID != 1 {
		t.Errorf("esperado ID 1, recebido %d", outputs[0].ID)
	}

	if outputs[0].Title != "Estudar MCP" {
		t.Errorf("esperado título %q, recebido %q", "Estudar MCP", outputs[0].Title)
	}

	if outputs[0].Completed {
		t.Error("primeira tarefa deveria estar pendente")
	}

	if outputs[1].ID != 2 {
		t.Errorf("esperado ID 2, recebido %d", outputs[1].ID)
	}

	if outputs[1].Title != "Estudar agentes" {
		t.Errorf("esperado título %q, recebido %q", "Estudar agentes", outputs[1].Title)
	}

	if !outputs[1].Completed {
		t.Error("segunda tarefa deveria estar concluída")
	}
}

func TestServiceListPropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("erro ao listar tarefas")

	repository := &repositoryFake{
		listFn: func(ctx context.Context) ([]domain.Task, error) {
			return nil, expectedErr
		},
	}

	service := NewService(repository)

	outputs, err := service.List(context.Background())

	if !errors.Is(err, expectedErr) {
		t.Fatalf("esperado erro %v, recebido %v", expectedErr, err)
	}

	if outputs != nil {
		t.Errorf("esperado resultado nil, recebido %+v", outputs)
	}
}

func TestServiceComplete(t *testing.T) {
	repository := &repositoryFake{
		completeFn: func(ctx context.Context, id int) (domain.Task, error) {
			task, err := domain.New(id, "Estudar MCP")
			if err != nil {
				return domain.Task{}, err
			}

			task.Complete()

			return task, nil
		},
	}

	service := NewService(repository)

	output, err := service.Complete(context.Background(), 1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if output.ID != 1 {
		t.Errorf("esperado ID 1, recebido %d", output.ID)
	}

	if output.Title != "Estudar MCP" {
		t.Errorf("esperado título %q, recebido %q", "Estudar MCP", output.Title)
	}

	if !output.Completed {
		t.Error("tarefa deveria estar concluída")
	}
}

func TestServiceCompleteInvalidID(t *testing.T) {
	repositoryCalled := false

	repository := &repositoryFake{
		completeFn: func(ctx context.Context, id int) (domain.Task, error) {
			repositoryCalled = true
			return domain.Task{}, nil
		},
	}

	service := NewService(repository)

	output, err := service.Complete(context.Background(), 0)

	if !errors.Is(err, domain.ErrInvalidID) {
		t.Fatalf("esperado ErrInvalidID, recebido %v", err)
	}

	if repositoryCalled {
		t.Error("repository não deveria ser chamado para ID inválido")
	}

	if output != (TaskOutput{}) {
		t.Errorf("esperado TaskOutput vazio, recebido %+v", output)
	}
}

func TestServiceCompleteNotFound(t *testing.T) {
	repository := &repositoryFake{
		completeFn: func(ctx context.Context, id int) (domain.Task, error) {
			return domain.Task{}, domain.ErrNotFound
		},
	}

	service := NewService(repository)

	output, err := service.Complete(context.Background(), 99)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("esperado ErrNotFound, recebido %v", err)
	}

	if output != (TaskOutput{}) {
		t.Errorf("esperado TaskOutput vazio, recebido %+v", output)
	}
}

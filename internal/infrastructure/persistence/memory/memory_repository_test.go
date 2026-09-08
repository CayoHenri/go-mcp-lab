package task

import (
	"context"
	"errors"
	"testing"

	domain "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

func TestMemoryRepositoryCreate(t *testing.T) {
	repository := NewMemoryRepository()

	created, err := repository.Create(context.Background(), "Estudar MCP")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if created.ID() != 1 {
		t.Errorf("esperado ID 1, recebido %d", created.ID())
	}

	if created.Title() != "Estudar MCP" {
		t.Errorf("esperado título Estudar MCP, recebido %s", created.Title())
	}

	if created.Completed() {
		t.Error("nova tarefa não deveria estar concluída")
	}
}

func TestMemoryRepositoryList(t *testing.T) {
	repository := NewMemoryRepository()
	ctx := context.Background()

	_, err := repository.Create(ctx, "Estudar MCP")
	if err != nil {
		t.Fatal(err)
	}

	_, err = repository.Create(ctx, "Estudar agentes")
	if err != nil {
		t.Fatal(err)
	}

	tasks, err := repository.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Fatalf("esperadas 2 tarefas, recebidas %d", len(tasks))
	}

	if tasks[0].ID() != 1 {
		t.Errorf("esperado primeiro ID 1, recebido %d", tasks[0].ID())
	}

	if tasks[1].ID() != 2 {
		t.Errorf("esperado segundo ID 2, recebido %d", tasks[1].ID())
	}
}

func TestMemoryRepositoryComplete(t *testing.T) {
	repository := NewMemoryRepository()
	ctx := context.Background()

	created, err := repository.Create(ctx, "Estudar MCP")
	if err != nil {
		t.Fatal(err)
	}

	completed, err := repository.Complete(ctx, created.ID())
	if err != nil {
		t.Fatal(err)
	}

	if !completed.Completed() {
		t.Error("tarefa deveria estar concluída")
	}
}

func TestMemoryRepositoryCompleteNotFound(t *testing.T) {
	repository := NewMemoryRepository()

	_, err := repository.Complete(context.Background(), 99)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("esperado ErrNotFound, recebido %v", err)
	}
}

func TestMemoryRepositoryFindByID(t *testing.T) {
	repository := NewMemoryRepository()
	ctx := context.Background()
	created, err := repository.Create(ctx, "Estudar MCP")
	if err != nil {
		t.Fatal(err)
	}

	found, err := repository.FindByID(ctx, created.ID())
	if err != nil {
		t.Fatal(err)
	}

	if found.ID() != created.ID() {
		t.Errorf("esperado ID %d, recebido %d", created.ID(), found.ID())
	}
}

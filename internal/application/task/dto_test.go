package task

import (
	"testing"

	domain "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

func TestNewTaskOutput(t *testing.T) {
	task, err := domain.New(10, "Estudar MCP")
	if err != nil {
		t.Fatalf("erro ao criar tarefa: %v", err)
	}

	task.Complete()

	output := NewTaskOutput(task)
	if output.ID != 10 {
		t.Errorf("esperado ID 10, recebido %d", output.ID)
	}

	if output.Title != "Estudar MCP" {
		t.Errorf("esperado título %q, recebido %q", "Estudar MCP", output.Title)
	}

	if !output.Completed {
		t.Error("esperado Completed true")
	}
}

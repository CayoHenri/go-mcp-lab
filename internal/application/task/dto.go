package task

import domain "github.com/CayoHenri/go-mcp-lab/internal/domain/task"

type TaskOutput struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func NewTaskOutput(task domain.Task) TaskOutput {
	return TaskOutput{
		ID:        task.ID(),
		Title:     task.Title(),
		Completed: task.Completed(),
	}
}

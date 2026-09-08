package postgres

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) domain.Repository {
	return &TaskRepository{db: db}
}

// Complete implements [task.Repository].
func (t *TaskRepository) Complete(ctx context.Context, id int) (domain.Task, error) {
	const query = `
		UPDATE tasks
		SET completed = TRUE
		WHERE id = $1
		RETURNING id, title, completed
	`

	var (
		taskID    int
		title     string
		completed bool
	)

	err := t.db.QueryRow(ctx, query, id).Scan(&taskID, &title, &completed)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Task{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.Task{}, fmt.Errorf("concluindo tarefa: %w", err)
	}

	return domain.Restore(taskID, title, completed), nil
}

// Create implements [task.Repository].
func (t *TaskRepository) Create(ctx context.Context, title string) (domain.Task, error) {
	const query = `
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id, title, completed
	`

	var (
		id        int
		taskTitle string
		completed bool
	)

	if err := t.db.QueryRow(ctx, query, title).Scan(&id, &taskTitle, &completed); err != nil {
		return domain.Task{}, fmt.Errorf("criando tarefa: %w", err)
	}

	return domain.Restore(id, taskTitle, completed), nil
}

// List implements [task.Repository].
func (t *TaskRepository) List(ctx context.Context) ([]domain.Task, error) {
	const query = `
		SELECT id, title, completed
		FROM tasks
		ORDER BY id
	`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("listando tarefas: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)

	for rows.Next() {
		var (
			id        int
			title     string
			completed bool
		)

		if err := rows.Scan(&id, &title, &completed); err != nil {
			return nil, fmt.Errorf("lendo tarefa: %w", err)
		}

		tasks = append(tasks, domain.Restore(
			id,
			title,
			completed,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("percorrendo tarefas: %w", err)
	}

	return tasks, nil
}

// FindByID implements [task.Repository].
func (t *TaskRepository) FindByID(ctx context.Context, id int) (domain.Task, error) {
	const query = `
		SELECT id, title, completed
		FROM tasks
		WHERE id = $1
	`
	var (
		taskID    int
		title     string
		completed bool
	)

	err := t.db.QueryRow(ctx, query, id).Scan(&taskID, &title, &completed)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Task{}, domain.ErrNotFound
	}
	return domain.Restore(taskID, title, completed), nil
}

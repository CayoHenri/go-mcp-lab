package task

import (
	"context"
	"sync"

	domain "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	nextID int
	tasks  []domain.Task
}

func NewMemoryRepository() domain.Repository {
	return &MemoryRepository{
		nextID: 1,
		tasks:  make([]domain.Task, 0),
	}
}

// Complete implements [task.Repository].
func (m *MemoryRepository) Complete(ctx context.Context, id int) (domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.tasks {
		if m.tasks[i].ID() != id {
			continue
		}

		m.tasks[i].Complete()
		return m.tasks[i], nil
	}

	return domain.Task{}, domain.ErrNotFound
}

// Create implements [task.Repository].
func (m *MemoryRepository) Create(ctx context.Context, title string) (domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, err := domain.New(m.nextID, title)
	if err != nil {
		return domain.Task{}, err
	}

	m.nextID++

	m.tasks = append(m.tasks, task)

	return task, nil
}

// List implements [task.Repository].
func (m *MemoryRepository) List(ctx context.Context) ([]domain.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]domain.Task, len(m.tasks))

	copy(tasks, m.tasks)

	return tasks, nil
}

// FindByID implements [task.Repository].
func (m *MemoryRepository) FindByID(ctx context.Context, id int) (domain.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, task := range m.tasks {
		if task.ID() == id {
			return task, nil
		}
	}
	return domain.Task{}, domain.ErrNotFound
}

// Delete implements [task.Repository].
func (m *MemoryRepository) Delete(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, task := range m.tasks {
		if task.ID() != id {
			continue
		}

		m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)

		return nil
	}

	return domain.ErrNotFound
}

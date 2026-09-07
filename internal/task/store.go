package task

import (
	"fmt"
	"sync"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Store struct {
	mu     sync.RWMutex
	nextID int
	tasks  []Task
}

func NewStore() *Store {
	return &Store{
		nextID: 1,
		tasks:  make([]Task, 0),
	}
}

func (s *Store) Create(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
	}

	s.nextID++

	s.tasks = append(s.tasks, task)

	return task
}

func (s *Store) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, len(s.tasks))

	copy(tasks, s.tasks)

	return tasks
}

func (s *Store) Complete(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID != id {
			continue
		}

		s.tasks[i].Completed = true

		return s.tasks[i], nil
	}

	return Task{}, fmt.Errorf("tarefa %d não encontrada", id)
}

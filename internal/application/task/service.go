package task

import (
	"context"

	domain "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, title string) (TaskOutput, error) {
	task, err := s.repository.Create(ctx, title)
	if err != nil {
		return TaskOutput{}, err
	}

	return NewTaskOutput(task), nil
}

func (s *Service) List(ctx context.Context) ([]TaskOutput, error) {
	tasks, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}
	outputs := make([]TaskOutput, len(tasks))
	for i, task := range tasks {
		outputs[i] = NewTaskOutput(task)
	}
	return outputs, nil
}

func (s *Service) Complete(ctx context.Context, id int) (TaskOutput, error) {
	if id <= 0 {
		return TaskOutput{}, domain.ErrInvalidID
	}

	task, err := s.repository.Complete(ctx, id)
	if err != nil {
		return TaskOutput{}, err
	}
	return NewTaskOutput(task), nil
}

func (s *Service) Summary(ctx context.Context) (SummaryOutput, error) {
	tasks, err := s.List(ctx)
	if err != nil {
		return SummaryOutput{}, err
	}

	output := SummaryOutput{
		Total: len(tasks),
		Tasks: tasks,
	}

	for _, task := range tasks {
		if task.Completed {
			output.Completed++
			continue
		}

		output.Pending++
	}

	return output, nil
}

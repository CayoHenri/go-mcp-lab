package task

import "context"

type Repository interface {
	Create(ctx context.Context, title string) (Task, error)
	List(ctx context.Context) ([]Task, error)
	Complete(ctx context.Context, id int) (Task, error)
	FindByID(ctx context.Context, id int) (Task, error)
	Delete(ctx context.Context, id int) error
}

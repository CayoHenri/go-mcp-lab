package task

import "strings"

type Task struct {
	id        int
	title     string
	completed bool
}

func New(id int, title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrTitleRequired
	}

	return Task{
		id:        id,
		title:     title,
		completed: false,
	}, nil
}

func Restore(id int, title string, completed bool) Task {
	return Task{
		id:        id,
		title:     title,
		completed: completed,
	}
}

func (t *Task) Complete() {
	t.completed = true
}

// Getters
func (t Task) ID() int {
	return t.id
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Completed() bool {
	return t.completed
}

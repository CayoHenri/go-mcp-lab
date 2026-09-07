package task

import "errors"

var (
	ErrNotFound      = errors.New("tarefa não encontrada")
	ErrTitleRequired = errors.New("o título da tarefa é obrigatório")
	ErrInvalidID     = errors.New("o identificador da tarefa deve ser maior que zero")
)

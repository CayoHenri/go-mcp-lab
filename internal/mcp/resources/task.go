package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	apptask "github.com/CayoHenri/go-mcp-lab/internal/application/task"
	domaintask "github.com/CayoHenri/go-mcp-lab/internal/domain/task"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const TaskTemplateURI = "tasks://{id}"

type TaskFinder interface {
	FindByID(ctx context.Context, id int) (apptask.TaskOutput, error)
}

func Task(service TaskFinder) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		id, err := taskIDFromURI(req.Params.URI)
		if err != nil {
			return nil, err
		}

		output, err := service.FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, domaintask.ErrNotFound) {
				return nil, fmt.Errorf("tarefa %d não encontrada", id)
			}

			return nil, fmt.Errorf("buscando tarefa %d: %w", id, err)
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("serializando tarefa: %w", err)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	}
}

// taskIDFromURI extracts the task ID from a given URI string.
// The expected format of the URI is "tasks://{id}" where {id} is a positive integer.
// If the URI is invalid or the ID is not a positive integer, an error is returned.
func taskIDFromURI(rawURI string) (int, error) {
	parsedURI, err := url.Parse(rawURI)
	if err != nil {
		return 0, fmt.Errorf("URI inválida: %w", err)
	}

	if parsedURI.Scheme != "tasks" {
		return 0, fmt.Errorf("scheme inválido: %s", parsedURI.Scheme)
	}

	idValue := parsedURI.Host
	if idValue == "" {
		idValue = strings.TrimPrefix(parsedURI.Path, "/")
	}

	id, err := strconv.Atoi(idValue)
	if err != nil {
		return 0, fmt.Errorf("ID da tarefa inválido: %s", idValue)
	}

	if id <= 0 {
		return 0, domaintask.ErrInvalidID
	}

	return id, nil
}

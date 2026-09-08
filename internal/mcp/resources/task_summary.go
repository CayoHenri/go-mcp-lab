package resources

import (
	"context"
	"encoding/json"
	"fmt"

	apptask "github.com/CayoHenri/go-mcp-lab/internal/application/task"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const TaskSummaryURI = "tasks://summary"

type TaskService interface {
	Summary(ctx context.Context) (apptask.SummaryOutput, error)
}

func TaskSummary(service TaskService) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		summary, err := service.Summary(ctx)
		if err != nil {
			return nil, fmt.Errorf("obtendo resumo das tarefas: %w", err)
		}

		data, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("serializando resumo das tarefas: %w", err)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      TaskSummaryURI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	}
}

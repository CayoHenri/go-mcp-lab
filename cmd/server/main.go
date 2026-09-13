package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	apptask "github.com/CayoHenri/go-mcp-lab/internal/application/task"
	"github.com/CayoHenri/go-mcp-lab/internal/config"
	postgresinfra "github.com/CayoHenri/go-mcp-lab/internal/infrastructure/persistence/postgres"
	mcpprompts "github.com/CayoHenri/go-mcp-lab/internal/mcp/prompts"
	mcpresources "github.com/CayoHenri/go-mcp-lab/internal/mcp/resources"
	"github.com/CayoHenri/go-mcp-lab/internal/mcp/tools"
)

func main() {
	boolPtr := func(value bool) *bool {
		return &value
	}

	ctx := context.Background()

	cfg, err := config.LoadServer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgresinfra.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := postgresinfra.EnsureSchema(ctx, db); err != nil {
		log.Fatal(err)
	}

	taskRepository := postgresinfra.NewTaskRepository(db)
	taskService := apptask.NewService(taskRepository)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "go-mcp-lab",
		Version: "v0.1.0",
	}, nil)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "greet",
			Description: "Cumprimenta uma pessoa pelo nome",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(false),
			},
		},
		tools.Greet,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "calculate",
			Description: "Realiza operações matemáticas básicas entre dois números",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(false),
			},
		},
		tools.Calculate,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "create_task",
			Description: "Cria uma nova tarefa",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				DestructiveHint: boolPtr(false),
				IdempotentHint:  false,
				OpenWorldHint:   boolPtr(false),
			},
		},
		tools.CreateTask(taskService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_tasks",
			Description: "Lista todas as tarefas",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				IdempotentHint:  true,
				DestructiveHint: boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		tools.ListTasks(taskService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "complete_task",
			Description: "Marca uma tarefa como concluída pelo identificador",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(false),
			},
		},
		tools.CompleteTask(taskService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "delete_task",
			Description: "Exclui permanentemente uma tarefa pelo identificador",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				DestructiveHint: boolPtr(true),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(false),
			},
		},
		tools.DeleteTask(taskService),
	)

	server.AddResource(
		&mcp.Resource{
			URI:         mcpresources.TaskSummaryURI,
			Name:        "task-summary",
			Title:       "Resumo das tarefas",
			Description: "Resumo das tarefas existentes, incluindo totais pendentes e concluídos",
			MIMEType:    "application/json",
		},
		mcpresources.TaskSummary(taskService),
	)

	server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: mcpresources.TaskTemplateURI,
			Name:        "task-by-id",
			Title:       "Tarefa por ID",
			Description: "Retorna os dados de uma tarefa específica pelo identificador",
			MIMEType:    "application/json",
		},
		mcpresources.Task(taskService),
	)

	server.AddPrompt(
		&mcp.Prompt{
			Name:        mcpprompts.TaskReviewName,
			Title:       "Revisão das tarefas",
			Description: "Analisa o estado atual das tarefas",
		},
		mcpprompts.TaskReview,
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

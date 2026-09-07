package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	apptask "github.com/CayoHenri/go-mcp-lab/internal/application/task"
	"github.com/CayoHenri/go-mcp-lab/internal/config"
	postgresinfra "github.com/CayoHenri/go-mcp-lab/internal/infrastructure/persistence/postgres"
	"github.com/CayoHenri/go-mcp-lab/internal/mcp/tools"
)

func main() {
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
		},
		tools.Greet,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "calculate",
			Description: "Realiza operações matemáticas básicas entre dois números",
		},
		tools.Calculate,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "create_task",
			Description: "Cria uma nova tarefa",
		},
		tools.CreateTask(taskService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_tasks",
			Description: "Lista todas as tarefas existentes",
		},
		tools.ListTasks(taskService),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "complete_task",
			Description: "Marca uma tarefa como concluída pelo identificador",
		},
		tools.CompleteTask(taskService),
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

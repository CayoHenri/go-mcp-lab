package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/CayoHenri/go-mcp-lab/internal/mcp/tools"
	"github.com/CayoHenri/go-mcp-lab/internal/task"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name: "go-mcp-lab", Version: "v0.1.0",
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

	taskStore := task.NewStore()

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "create_task",
			Description: "Cria uma nova tarefa",
		},
		tools.CreateTask(taskStore),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_tasks",
			Description: "Lista todas as tarefas existentes",
		},
		tools.ListTasks(taskStore),
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "complete_task",
			Description: "Marca uma tarefa como concluída pelo identificador",
		},
		tools.CompleteTask(taskStore),
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

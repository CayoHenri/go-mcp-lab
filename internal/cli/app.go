package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/CayoHenri/go-mcp-lab/internal/agent"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
)

type App struct {
	mcp     *mcpclient.Client
	agent   *agent.Agent
	console *Console
}

func New(mcp *mcpclient.Client, agentClient *agent.Agent, console *Console) *App {
	return &App{
		mcp:     mcp,
		agent:   agentClient,
		console: console,
	}
}

func (a *App) Run(ctx context.Context) error {
	printHeader()

	for {
		input, err := a.console.ReadLine("> ")
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if input == "" {
			continue
		}

		shouldExit, err := a.handleInput(ctx, input)
		if err != nil {
			printError(err)
			continue
		}

		if shouldExit {
			break
		}
	}

	fmt.Println()
	fmt.Println("Até mais!")

	return nil
}

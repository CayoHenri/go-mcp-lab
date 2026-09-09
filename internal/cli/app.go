package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/CayoHenri/go-mcp-lab/internal/agent"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
)

type App struct {
	mcp   *mcpclient.Client
	agent *agent.Agent
}

func New(mcp *mcpclient.Client, agentClient *agent.Agent) *App {
	return &App{
		mcp:   mcp,
		agent: agentClient,
	}
}

func (a *App) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	printHeader()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
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

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Até mais!")

	return nil
}

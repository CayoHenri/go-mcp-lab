package agent

import (
	"context"
	"fmt"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
	mcpclient "github.com/CayoHenri/go-mcp-lab/internal/mcp/client"
)

type Agent struct {
	mcp      *mcpclient.Client
	llm      *llm.Client
	session  *Session
	approver Approver
	policy   *PermissionPolicy
}

func New(mcp *mcpclient.Client, llmClient *llm.Client, approver Approver) *Agent {
	return &Agent{
		mcp:      mcp,
		llm:      llmClient,
		session:  NewSession(),
		approver: approver,
		policy:   NewPermissionPolicy(),
	}
}

func (a *Agent) Run(ctx context.Context, input string) (string, error) {
	return a.runLoop(ctx, input)
}

func (a *Agent) RunPrompt(ctx context.Context, promptName string, arguments map[string]string) (string, error) {
	prompt, err := a.mcp.GetPrompt(ctx, promptName, arguments)
	if err != nil {
		return "", fmt.Errorf("obtendo MCP prompt %s: %w", promptName, err)
	}

	input, err := promptToText(prompt)
	if err != nil {
		return "", fmt.Errorf("convertendo MCP prompt %s: %w", promptName, err)
	}

	return a.runLoop(ctx, input)
}

func (a *Agent) Reset() {
	a.session.Reset()
}

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	llm "github.com/CayoHenri/go-mcp-lab/internal/llm/openai"
)

type Approver struct {
	console *Console
}

func NewApprover(console *Console) *Approver {
	return &Approver{
		console: console,
	}
}

func (a *Approver) Approve(ctx context.Context, call llm.FunctionCall) (bool, error) {
	arguments, err := json.MarshalIndent(call.Arguments, "", "  ")
	if err != nil {
		return false, fmt.Errorf("serializando argumentos: %w", err)
	}

	fmt.Println()
	fmt.Println("O agente deseja executar uma ação:")

	fmt.Println()
	fmt.Printf("Função: %s\n", call.Name)

	fmt.Printf("Argumentos:\n%s\n", arguments)

	fmt.Println()

	answer, err := a.console.ReadLine("Autorizar? [s/N]: ")
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(strings.TrimSpace(answer))

	return answer == "s" || answer == "sim" || answer == "y" || answer == "yes", nil
}

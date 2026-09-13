package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CayoHenri/go-mcp-lab/internal/agent"
)

type Approver struct {
	console *Console
}

func NewApprover(console *Console) *Approver {
	return &Approver{
		console: console,
	}
}

func (a *Approver) Approve(ctx context.Context, request agent.ApprovalRequest) (bool, error) {
	call := request.Call

	arguments, err := json.MarshalIndent(call.Arguments, "", "  ")
	if err != nil {
		return false, err
	}

	fmt.Println()
	fmt.Println("O agente deseja executar uma ação:")

	fmt.Println()

	fmt.Printf("Função: %s\n", call.Name)

	fmt.Printf("Permissão: %s\n", request.PermissionLevel)

	fmt.Printf("Argumentos:\n%s\n", arguments)

	if request.PermissionLevel == agent.PermissionDestructive {
		fmt.Println()
		fmt.Println("ATENÇÃO: esta ação é destrutiva e pode remover dados permanentemente.")
	}

	fmt.Println()

	answer, err := a.console.ReadLine("Autorizar? [s/N]: ")
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(strings.TrimSpace(answer))

	return answer == "s" || answer == "sim" || answer == "y" || answer == "yes", nil
}

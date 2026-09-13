package cli

import (
	"context"
	"fmt"
	"strings"
)

func (a *App) handleInput(ctx context.Context, input string) (bool, error) {
	if !strings.HasPrefix(input, "/") {
		return false, a.runAgent(ctx, input)
	}

	parts := strings.Fields(input)

	command := strings.ToLower(parts[0])

	switch command {
	case "/exit":
		return true, nil

	case "/help":
		printHelp()
		return false, nil

	case "/tools":
		return false, a.printTools(ctx)

	case "/resources":
		return false, a.printResources(ctx)

	case "/prompts":
		return false, a.printPrompts(ctx)

	case "/prompt":
		return false, a.runPrompt(ctx, parts)
	case "/reset":
		a.agent.Reset()

		fmt.Println()
		fmt.Println("Conversa reiniciada.")
		fmt.Println()

		return false, nil
	default:
		fmt.Printf("Comando desconhecido: %s\n", command)

		fmt.Println("Digite /help para ver os comandos disponíveis.")

		fmt.Println()

		return false, nil
	}
}

func (a *App) runAgent(ctx context.Context, input string) error {
	response, err := a.agent.Run(ctx, input)
	if err != nil {
		return err
	}

	printResponse(response)

	return nil
}

func (a *App) runPrompt(ctx context.Context, parts []string) error {
	if len(parts) < 2 {
		fmt.Println()
		fmt.Println("Uso: /prompt <nome>")
		fmt.Println()

		return nil
	}

	promptName := parts[1]

	response, err := a.agent.RunPrompt(ctx, promptName, nil)
	if err != nil {
		return err
	}

	printResponse(response)

	return nil
}

func (a *App) printTools(ctx context.Context) error {
	result, err := a.mcp.ListTools(ctx)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("MCP Tools:")

	if len(result.Tools) == 0 {
		fmt.Println("Nenhuma tool disponível.")
		fmt.Println()

		return nil
	}

	for _, tool := range result.Tools {
		fmt.Printf("- %s", tool.Name)

		if tool.Description != "" {
			fmt.Printf(": %s", tool.Description)
		}

		fmt.Println()

		if tool.Annotations != nil {
			fmt.Printf("  readOnly: %t\n", tool.Annotations.ReadOnlyHint)
			fmt.Printf("  destructive: %t\n", destructiveValue(tool.Annotations.DestructiveHint))
			fmt.Printf("  idempotent: %t\n", tool.Annotations.IdempotentHint)
		}
	}

	fmt.Println()

	return nil
}

func (a *App) printResources(ctx context.Context) error {
	resources, err := a.mcp.ListResources(ctx)
	if err != nil {
		return err
	}

	templates, err := a.mcp.ListResourceTemplates(ctx)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("MCP Resources:")

	if len(resources) == 0 {
		fmt.Println("Nenhum resource disponível.")
	}

	for _, resource := range resources {

		fmt.Printf("- %s", resource.URI)

		if resource.Description != "" {
			fmt.Printf(": %s", resource.Description)
		}

		fmt.Println()
	}

	fmt.Println()
	fmt.Println("MCP Resource Templates:")

	if len(templates) == 0 {
		fmt.Println("Nenhum resource template disponível.")
	}

	for _, template := range templates {

		fmt.Printf("- %s", template.URITemplate)

		if template.Description != "" {
			fmt.Printf(": %s", template.Description)
		}

		fmt.Println()
	}

	fmt.Println()

	return nil
}

func (a *App) printPrompts(ctx context.Context) error {
	prompts, err := a.mcp.ListPrompts(ctx)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("MCP Prompts:")

	if len(prompts) == 0 {
		fmt.Println("Nenhum prompt disponível.")

		fmt.Println()

		return nil
	}

	for _, prompt := range prompts {

		fmt.Printf("- %s", prompt.Name)

		if prompt.Description != "" {
			fmt.Printf(": %s", prompt.Description)
		}

		fmt.Println()
	}

	fmt.Println()

	return nil
}

func destructiveValue(value *bool) bool {
	if value == nil {
		return true
	}

	return *value
}

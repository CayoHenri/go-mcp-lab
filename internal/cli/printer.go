package cli

import "fmt"

func printHeader() {
	fmt.Println()
	fmt.Println("Go MCP Lab")
	fmt.Println("==========")
	fmt.Println()

	fmt.Println("Converse com o agente ou digite /help para ver os comandos.")
	fmt.Println()
}

func printResponse(response string) {
	fmt.Println()
	fmt.Println(response)
	fmt.Println()
}

func printError(err error) {
	fmt.Println()

	fmt.Printf("Erro: %v\n", err)

	fmt.Println()
}

func printHelp() {
	fmt.Println()
	fmt.Println("Comandos disponíveis:")

	fmt.Println()

	fmt.Println("  /help              Exibe esta ajuda")

	fmt.Println("  /tools             Lista MCP Tools")

	fmt.Println("  /resources         Lista MCP Resources")

	fmt.Println("  /prompts           Lista MCP Prompts")

	fmt.Println("  /prompt <nome>     Executa um MCP Prompt")

	fmt.Println("  /exit              Encerra o cliente")

	fmt.Println("  /reset             Inicia uma nova conversa")

	fmt.Println()
}

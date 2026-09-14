# 11. Resumo e Checklist

[← Voltar ao README](../README.md)

## 11.1 Resumo para revisão rápida

Se passar alguns meses sem trabalhar no projeto, releia isto:

| Conceito | Definição |
|---|---|
| **LLM** | interpreta e decide |
| **Agent** | coordena o ciclo de decisão e execução |
| **MCP Client** | conversa com MCP Servers |
| **MCP Server** | publica capacidades e contexto |
| **Tool** | executa uma capacidade |
| **Resource** | fornece informação/contexto |
| **Prompt** | fornece uma orientação reutilizável |
| **Function Call** | pedido estruturado produzido pela LLM |
| **Agent Loop** | continua executando e consultando a LLM até obter resposta final |
| **Session** | mantém contexto entre mensagens |
| **HITL** | pede decisão humana antes de ações sensíveis |
| **Tool Annotations** | fornecem sinais sobre comportamento da Tool |
| **Permission Policy** | decide como uma capacidade pode ser executada |
| **Deny by default** | capacidade desconhecida não ganha permissão automaticamente |

A frase principal para guardar é:

> **A LLM decide o que fazer, o Agent coordena como continuar, e MCP padroniza como capacidades e contexto são disponibilizados.**

---

## 11.2 Checklist para retomar o projeto

- [ ] Subir ambiente
- [ ] Rodar `go fmt`
- [ ] Rodar `go vet`
- [ ] Rodar `go test`
- [ ] Executar `cmd/client`
- [ ] Ver `/tools`
- [ ] Ver `/resources`
- [ ] Ver `/prompts`
- [ ] Testar READ
- [ ] Testar WRITE + approval
- [ ] Testar multi-turn
- [ ] Testar DESTRUCTIVE + approval
- [ ] Testar `/prompt task-review`
- [ ] Testar `/reset`

Se tudo funcionar, o próximo estudo é **Multiple MCP Servers** — veja [`10-roadmap-e-proximos-passos.md`](10-roadmap-e-proximos-passos.md).

---

[← Roadmap e Próximos Passos](10-roadmap-e-proximos-passos.md) · [Voltar ao README](../README.md)

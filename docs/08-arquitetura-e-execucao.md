# 8. Arquitetura Atual e Como Rodar

[← Voltar ao README](../README.md)

## 8.1 O Agent que construímos até aqui

A arquitetura conceitual atual é:

```mermaid
flowchart TD
    U[Usuário] --> CLI
    CLI --> Agent
    Agent --> Session
    Agent --> LLM
    LLM --> FC[Function Calls]
    FC --> Auth[Authorization]
    Auth --> READ
    Auth --> WD["WRITE / DESTRUCTIVE"]
    WD --> HITL
    READ --> FE[Function Executor]
    HITL --> FE
    FE --> MT[MCP Tool]
    FE --> RR[read_resource]
    MT --> Client[MCP Client]
    RR --> Client
    Client --> Server[MCP Server]
    Server --> Tools
    Server --> Resources
    Server --> Prompts
```

---

## 8.2 Como rodar o laboratório

### Pré-requisitos

```text
Go
Docker
Git
OpenAI API Key
```

Entre no projeto:

```powershell
cd C:\desenvolvimento\go-mcp-lab
```

Configure o `.env`:

```env
OPENAI_API_KEY=sua-chave-aqui
DATABASE_URL=postgres://postgres:postgres@localhost:5432/mcp_lab?sslmode=disable
```

Suba os serviços necessários:

```powershell
docker compose up -d
```

Valide o projeto:

```powershell
go fmt ./...
go vet ./...
go test ./...
```

Execute:

```powershell
go run ./cmd/client
```

---

## 8.3 Comandos da CLI

| Comando | Efeito |
|---|---|
| `/help` | Mostra ajuda |
| `/tools` | Lista MCP Tools descobertas |
| `/resources` | Lista Resources |
| `/prompts` | Lista Prompts |
| `/prompt task-review` | Executa o Prompt MCP |
| `/reset` | Reinicia o contexto conversacional |
| `/exit` | Encerra |

---

## 8.4 Roteiro rápido para demonstrar tudo

Ao retomar o projeto, execute em sequência:

```text
/tools
/resources
/prompts
```

Agora teste uma consulta:

```text
Liste minhas tarefas.
```

Teste uma escrita:

```text
Crie uma tarefa chamada Revisar MCP.
```

Teste multi-turn:

```text
Agora conclua ela.
```

Teste ação destrutiva:

```text
Agora exclua ela.
```

Teste Prompt:

```text
/prompt task-review
```

E finalmente:

```text
/reset
```

Esse pequeno roteiro exercita quase todos os conceitos estudados.

---

[← Human-in-the-loop e Permissões](07-human-in-the-loop-e-permissoes.md) · [Próximo: Boas Práticas →](09-boas-praticas.md)

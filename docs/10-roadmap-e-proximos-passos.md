# 10. Roadmap e Próximos Passos

[← Voltar ao README](../README.md)

## 10.1 Como pensar em MCP durante novos projetos

Ao encontrar uma funcionalidade, pergunte:

```mermaid
flowchart TD
    Q{"Que tipo de\nfuncionalidade é?"}
    Q -->|É uma ação?| Tool["Tool\ncreate_issue, send_email,\ndeploy_application"]
    Q -->|É informação/contexto?| Resource["Resource\nrepository://readme,\nproject://architecture, customer://123"]
    Q -->|É orientação reutilizável?| Prompt["Prompt\nreview-pull-request,\nanalyze-incident, prepare-release"]
```

---

## 10.2 Quando uma aplicação começa a parecer um Agent?

Somente utilizar uma LLM não significa necessariamente ter um Agent.

Exemplo simples:

```mermaid
flowchart LR
    User --> LLM --> Text
```

Com ferramentas:

```mermaid
flowchart LR
    User --> LLM --> Tool --> Text
```

Com Agent Loop:

```mermaid
flowchart LR
    User --> LLM1[LLM] --> Tool1[Tool] --> LLM2[LLM] --> Tool2[Tool] --> LLM3[LLM] --> Final
```

Agora existe um processo iterativo de decisão e execução. Adicionando `state`, `permissions`, `human approval` e `external capabilities`, temos uma arquitetura agentic cada vez mais completa.

---

## 10.3 Autonomia controlada

Um Agent útil não precisa ter autonomia irrestrita. Nosso modelo é:

```mermaid
flowchart TD
    Agent["Agent pode decidir\nqual capacidade usar"] -->|MAS| Sens["ações sensíveis\npodem exigir aprovação"]
```

Isso cria **Autonomia + Controle**.

Human-in-the-loop não reduz o conceito de Agent — ele define limites para a autonomia.

---

## 10.4 Próximo grande passo

O projeto foi pausado neste ponto. A próxima evolução recomendada é: **múltiplos MCP Servers.**

Hoje:

```mermaid
flowchart LR
    Agent --> Task[Task MCP Server]
```

Depois:

```mermaid
flowchart TD
    Agent --> Reg[MCP Registry]
    Reg --> Task[Task MCP]
    Reg --> GH[GitHub MCP]
```

Isso permitirá estudar:

- múltiplas sessões MCP;
- discovery agregado;
- roteamento de Tools;
- origem de cada Tool;
- colisão de nomes;
- trust por servidor;
- autorização baseada em origem;
- Tool Annotations de servidores diferentes.

---

## 10.5 Evolução futura de segurança

Com vários servidores:

```mermaid
flowchart TD
    Tool --> Q1{"Qual servidor publicou?"}
    Q1 --> Q2{"Eu confio nesse servidor?"}
    Q2 --> Q3{"Quais annotations a Tool possui?"}
    Q3 --> Q4{"Qual é minha policy local?"}
    Q4 --> Q5{"Precisa de aprovação?"}
```

Arquitetura:

```mermaid
flowchart LR
    SI[Server Identity] --> AD[Authorization Decision]
    ST[Server Trust] --> AD
    TA[Tool Annotations] --> AD
    LP[Local Policy] --> AD
```

Esse será um passo importante para sair de um laboratório simples para um Agent mais próximo de aplicações reais.

---

## 10.6 Outros assuntos para estudar depois

Após múltiplos MCP Servers:

1. Server Registry
2. Server Trust
3. Streaming
4. Elicitation
5. Sampling
6. Observability
7. Agent testing
8. Web UI
9. Agent memory
10. MCP remoto

Cada tema deve ser introduzido isoladamente para entender claramente sua responsabilidade.

---

## 10.7 Mapa do aprendizado

```mermaid
flowchart TD
    A[MCP Server] --> B[MCP Client] --> C[Tools] --> D[LLM Integration]
    D --> E[Function Calling] --> F[Agent Loop] --> G[Multiple Function Calls]
    G --> H[Resources] --> I[Resource Templates] --> J[Prompts]
    J --> K[Interactive CLI] --> L[Multi-turn] --> M[Human-in-the-loop]
    M --> N[Permission Policy] --> O[Destructive Actions] --> P[Tool Annotations]
    P --> Pause["⏸ PAUSA ATUAL"]
    Pause --> Q[Multiple MCP Servers] --> R["Trust & Security"] --> S["Advanced MCP / Agents"]
```

---

[← Boas Práticas](09-boas-praticas.md) · [Próximo: Resumo e Checklist →](11-resumo-e-checklist.md)

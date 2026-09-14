# syntax=docker/dockerfile:1

# ---------- Build stage ----------
FROM golang:1.25.5-alpine AS builder

ENV GOTOOLCHAIN=local

WORKDIR /app

# Cache de dependências: copia só go.mod/go.sum primeiro
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila os DOIS binários: client e server.
# O client, em runtime, não usa mais "go run" — ele invoca o binário do
# server diretamente (via MCP_SERVER_BIN, ver internal/client/client.go).
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /app/bin/client ./cmd/client && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /app/bin/server ./cmd/server

# ---------- Runtime stage ----------
# Agora que não depende mais do toolchain do Go em runtime, voltamos pra
# distroless: sem shell, sem package manager, só os binários + certificados
# TLS. Imagem final na casa de 15-25 MB em vez de ~900 MB.
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /app

COPY --from=builder /app/bin/client /app/client
COPY --from=builder /app/bin/server /app/server

# Diz pro client usar o binário compilado em vez de "go run ./cmd/server"
ENV MCP_SERVER_BIN=/app/server

USER nonroot:nonroot

# CLI interativa: mantenha stdin/tty abertos ao rodar o container
ENTRYPOINT ["/app/client"]
# Multi-stage build para aplicação Go leve e segura
FROM golang:alpine AS builder

WORKDIR /app

# Instala certificados e git se necessário
RUN apk add --no-cache ca-certificates git

# Cache de dependências do Go
COPY go.mod go.sum ./
RUN go mod download

# Copia código fonte
COPY . .

# Compila binário estático e otimizado
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /dynamicqr ./cmd/server/main.go

# Imagem final de execução minimalista
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copia o executável e os arquivos estáticos do frontend
COPY --from=builder /dynamicqr /app/dynamicqr
COPY --from=builder /app/web/static /app/web/static

EXPOSE 8080

ENTRYPOINT ["/app/dynamicqr"]

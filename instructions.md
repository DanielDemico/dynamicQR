# Instructions — DynamicQR

Este arquivo é a referência operacional para configurar, executar, testar e operar o sistema **DynamicQR**.

---

## 1. Pré-requisitos

Para executar o projeto são necessários:
- **Docker** (versão 24+ com Docker Compose v2)
- **Go** (versão 1.22 ou superior, caso queira rodar localmente sem Docker)
- **Git**

---

## 2. Variáveis de Ambiente

Copie o arquivo de exemplo para criar a sua configuração:
```bash
cp .env.example .env
```

### Principais variáveis:
- `PORT`: Porta HTTP onde o servidor escuta (padrão: `8080`).
- `BASE_URL`: URL base do servidor utilizada para compor as URLs curtas e os QR codes.
  - Para testes locais imediatos: `http://localhost:8080`
  - Para produção: `https://dynamicqr.danielcezar.com.br`
- `MONGO_URI`: String de conexão com o MongoDB (`mongodb://mongo:27017` no Docker, `mongodb://localhost:27017` fora do Docker).
- `MONGO_DB`: Nome do banco no MongoDB (`dynamicqr_db`).

---

## 3. Como Executar

### 3.1 Via Docker Compose (Recomendado)
Para subir o banco de dados MongoDB e o servidor Go compilado:
```bash
docker compose up --build -d
```

Para visualizar os logs da aplicação:
```bash
docker compose logs -f app
```

Para parar o ambiente mantendo os dados persistidos:
```bash
docker compose down
```

Para parar e remover volumes (limpeza total da base de dados):
```bash
docker compose down -v
```

### 3.2 Execução Local em Desenvolvimento (Sem Docker para o App Go)
Caso queira executar apenas o MongoDB no Docker e o binário Go diretamente na máquina:
```bash
# Subir apenas o MongoDB
docker compose up -d mongo

# Executar a aplicação Go
go run cmd/server/main.go
```

---

## 4. Portas e Serviços

| Serviço | Porta do Host | Porta do Container | Descrição |
|---|---|---|---|
| `app` | `8080` | `8080` | Servidor Web Go (API REST + Frontend + Redirecionador) |
| `mongo` | `27017` | `27017` | Banco de Dados NoSQL MongoDB 7 com volume persistente |

---

## 5. Como Executar os Testes

### 5.1 Testes Unitários
```bash
go test -v ./...
```

### 5.2 Testes End-to-End (E2E)
Com o sistema em execução (na porta 8080):

**No Windows (PowerShell):**
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/e2e/spec-001-qr-dinamico.ps1
```

**No Linux / macOS (Bash):**
```bash
./scripts/e2e/spec-001-qr-dinamico.sh
```

---

## 6. Estrutura do Projeto

```
DynamicQR/
├── _SPECS/                          # Especificações formais (Spec-Driven Development)
│   ├── INDEX.md                     # Índice central das specs
│   └── qr-code-dinamico/            # Spec SPEC-001
│       ├── SPEC.md
│       └── contract.md
├── cmd/
│   └── server/
│       └── main.go                  # Ponto de entrada do executável Go
├── internal/
│   ├── config/                      # Leitura de variáveis de ambiente
│   ├── database/                    # Conexão MongoDB e criação de índices únicos
│   ├── domain/                      # Entidades (QRCode com campo Base32)
│   ├── handler/                     # Rotas REST e redirecionamento público
│   ├── repository/                  # Operações no MongoDB (Insert, Find, Update, Inc)
│   └── service/                     # Geração de QR Code e encode/decode Base32
├── web/
│   └── static/                      # Frontend Vanilla (HTML5 / CSS / JS)
├── scripts/
│   └── e2e/                         # Scripts de validação E2E
├── Dockerfile                       # Multi-stage build da aplicação Go
├── docker-compose.yml               # Orquestração do app + MongoDB
├── instructions.md                  # Este manual operacional
├── AGENTS.MD                        # Diretrizes do projeto e SDD
└── README.MD
```

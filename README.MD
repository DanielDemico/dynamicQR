# DynamicQR ⚡

Sistema de alta performance para criação, gestão e redirecionamento de **QR Codes Dinâmicos**, desenvolvido em **Go**, persistido em **MongoDB** (com codificação **Base32** das imagens) e conteinerizado via **Docker Compose**.

Interface de controle moderna desenvolvida em **HTML5 + Vanilla CSS + Vanilla JS**.

---

## 🚀 Inicialização Rápida

Para subir o banco de dados e a aplicação web:

```bash
docker compose up --build -d
```

Acesse o painel em seu navegador:
👉 **[http://localhost:8080](http://localhost:8080)**

---

## 🧪 Testes e Validação (Spec-Driven Development)

### Testes Unitários
```bash
go test -v ./...
```

### Testes End-to-End (E2E)
Com o sistema em execução:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/e2e/spec-001-qr-dinamico.ps1
```

---

## 📖 Documentação Completa

- [Manual Operacional e de Configuração](instructions.md)
- [Diretrizes de Spec-Driven Development](AGENTS.MD)
- [Índice de Especificações](_SPECS/INDEX.md)
- [SPEC-001: Gerador de QR Code Dinâmico](_SPECS/qr-code-dinamico/SPEC.md)
- [Contrato Técnico da API e MongoDB](_SPECS/qr-code-dinamico/contract.md)
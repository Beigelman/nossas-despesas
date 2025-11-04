# Nossas Despesas - Monorepo

Este monorepo contém os serviços do sistema Nossas Despesas para organizar e dividir despesas em grupo:

- **Backend API** (Go): API principal com autenticação, grupos, despesas e receitas
- **ML Service** (Python): Serviço de machine learning para classificação automática de categorias de despesas

## Funcionalidades Principais

- **Autenticação**: login por credenciais ou Google e geração/renovação de tokens.
- **Usuários**: consulta do usuário logado.
- **Grupos**: criação de grupos, envio de convites e consulta de saldo do grupo.
- **Categorias**: cadastro de categorias e grupos de categorias de despesas.
- **Despesas**: criação, atualização, exclusão e agendamento de despesas, além de geração de relatórios por período e categoria.
- **Receitas**: registro de receitas e consulta mensal.

## Estrutura do Monorepo

```
/
├── backend/              # Backend API (Go)
│   ├── internal/         # Código interno do backend
│   ├── database/         # Migrações e scripts de BD
│   ├── templates/        # Templates de email
│   └── scripts/          # Scripts utilitários
├── machine_learn/        # Serviço ML (Python)
│   ├── src/              # Código fonte da API ML
│   ├── training/         # Scripts de treinamento
│   └── models/           # Modelos treinados
└── .github/workflows/    # CI/CD separado por serviço
```

### Backend API (Go)

O backend utiliza o micro framework **Eon** (pasta `backend/internal/pkg/eon`) para gerenciar ciclos de vida e injeção de dependências. Cada domínio possui um módulo próprio sob `backend/internal/modules`, contendo:

- `controller/` – handlers HTTP definidos com o framework Fiber.
- `usecase/` – regras de negócio.
- `postgres/` – repositórios com implementação para Postgres.
- `module/` – arquivo que registra o módulo no Eon para bootstrapping.

Módulos comuns (servidor HTTP, configuração, banco de dados) ficam em `backend/internal/pkg`. Há ainda um módulo `backend/internal/shared` que provê clientes de infraestrutura (JWT, envio de e-mail, pub/sub) e middlewares.

As migrações do banco estão em `backend/database/migrations` e podem ser executadas pelo script `backend/database/migrate.sh`.

### ML Service (Python)

O serviço de ML utiliza **FastAPI** e **Poetry** para classificação automática de categorias de despesas. Inclui:

- API REST para predição em tempo real e batch
- Pipeline de treinamento de modelos
- Modelos persistidos com joblib/pickle

## Execução Local

### Backend API (Go)

1. Inicie o Postgres via Docker Compose:
   ```bash
   docker compose up db -d
   ```
2. Aplique as migrações:
   ```bash
   cd backend
   ./database/migrate.sh up ./database/migrations
   ```
3. Inicie a aplicação:
   ```bash
   cd backend
   ENV=development go run main.go
   ```

### ML Service (Python)

1. Instale as dependências:
   ```bash
   cd machine_learn
   poetry install
   ```
2. Inicie o serviço:
   ```bash
   cd machine_learn
   poetry run uvicorn src.main:app --host 0.0.0.0 --port 8000 --reload
   ```

Cada serviço possui seu próprio `Makefile` com atalhos para as tarefas comuns.

## Testes e Qualidade de Código

O repositório possui workflows do GitHub Actions separados para cada serviço (diretório `.github/workflows`).

### Backend (Go)
```bash
cd backend
go test ./...
golangci-lint run
```

### ML Service (Python)
```bash
cd machine_learn
poetry run pytest
poetry run flake8 src/
poetry run mypy src/
```

## Deploy

O monorepo possui workflows de deploy independentes que são acionados apenas quando há mudanças no respectivo serviço:

- **Backend**: `backend-deploy.yml` faz deploy da API Go no Cloud Run (`nossas-despesas-be`)
- **ML Service**: `ml-deploy.yml` faz deploy do serviço Python no Cloud Run (`nossas-despesas-ml`)

Cada serviço possui versionamento semântico independente com prefixos `backend-v` e `ml-v`.


# Nossas Despesas - Monorepo

Este monorepo contém os serviços do sistema Nossas Despesas para organizar e dividir despesas em grupo:

- **Backend API** (Go): API principal com autenticação, grupos, despesas e receitas
- **ML Service** (Python): Serviço de machine learning para classificação automática de categorias de despesas
- **Web App** (Next.js): Aplicação web frontend construída com Next.js 15, React 19 e TypeScript

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
├── web/                  # Web App (Next.js)
│   ├── src/              # Código fonte da aplicação
│   ├── public/           # Arquivos estáticos
│   └── vercel.json       # Configuração do Vercel
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

### Web App (Next.js)

A aplicação web utiliza **Next.js 15**, **React 19**, **TypeScript** e **Tailwind CSS**. Inclui:

- Interface moderna e responsiva com componentes Radix UI
- Autenticação via NextAuth com suporte a Google OAuth e credenciais
- PWA (Progressive Web App) com suporte offline
- Integração com backend API via axios
- Gerenciamento de estado com Zustand e React Query

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

### Web App (Next.js)

1. Instale as dependências:
   ```bash
   cd web
   pnpm install
   ```

2. Configure as variáveis de ambiente (veja seção [Variáveis de Ambiente](#variáveis-de-ambiente)):
   ```bash
   cd web
   cp .env.example .env.local
   # Edite .env.local com suas configurações
   ```

3. Inicie o servidor de desenvolvimento:
   ```bash
   cd web
   pnpm dev
   ```

A aplicação estará disponível em `http://localhost:3000`.

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

### Web App (Next.js)
```bash
cd web
pnpm lint          # Executa ESLint
pnpm exec tsc --noEmit  # Verifica tipos TypeScript
# TODO: Adicionar testes quando implementados
# pnpm test
```

## Deploy

O monorepo possui workflows de deploy independentes que são acionados apenas quando há mudanças no respectivo serviço:

- **Backend**: `backend-deploy.yml` faz deploy da API Go no Cloud Run (`nossas-despesas-be`)
- **ML Service**: `ml-deploy.yml` faz deploy do serviço Python no Cloud Run (`nossas-despesas-ml`)
- **Web App**: `web-deploy.yml` faz deploy da aplicação Next.js no Vercel

Cada serviço possui versionamento semântico independente com prefixos `backend-v` e `ml-v`. O Vercel gerencia versões automaticamente para o web app.

## Variáveis de Ambiente

### Web App (Next.js)

As seguintes variáveis de ambiente são necessárias para o funcionamento do app web:

#### Variáveis Públicas (NEXT_PUBLIC_*)
- `NEXT_PUBLIC_API_URL`: URL completa do backend API (ex: `https://api.example.com`)
- `NEXT_PUBLIC_BASE_URL`: URL base da aplicação web (ex: `https://app.example.com`)

#### Variáveis de Autenticação
- `GOOGLE_CLIENT_ID`: Client ID do Google OAuth (obtido no Google Cloud Console)
- `GOOGLE_CLIENT_SECRET`: Client Secret do Google OAuth (obtido no Google Cloud Console)
- `NEXTAUTH_SECRET`: Secret usado pelo NextAuth para criptografar tokens (gere com: `openssl rand -base64 32`)
- `NEXTAUTH_URL`: URL completa da aplicação (deve corresponder a `NEXT_PUBLIC_BASE_URL`)

#### Configuração no Vercel

Para configurar as variáveis no Vercel:

1. Acesse o projeto no [Vercel Dashboard](https://vercel.com)
2. Vá em **Settings** → **Environment Variables**
3. Adicione todas as variáveis listadas acima
4. Certifique-se de que estão configuradas para **Production**, **Preview** e **Development** conforme necessário

#### Configuração no GitHub Actions

As seguintes secrets devem ser configuradas no GitHub:

- `VERCEL_TOKEN`: Token de autenticação do Vercel (obtido em [Vercel Settings → Tokens](https://vercel.com/account/tokens))
- `VERCEL_ORG_ID`: ID da organização no Vercel
- `VERCEL_PROJECT_ID`: ID do projeto no Vercel
- Todas as variáveis de ambiente listadas acima (para uso no build)

Para obter `VERCEL_ORG_ID` e `VERCEL_PROJECT_ID`:
```bash
# Instale o Vercel CLI
npm i -g vercel

# Faça login
vercel login

# Link o projeto
cd web
vercel link

# Os IDs estarão no arquivo .vercel/project.json
```


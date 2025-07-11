# Nossas Despesas API

Nossas Despesas é uma API escrita em Go para organizar e dividir despesas em grupo. O projeto oferece autenticação por JWT, cadastro de usuários, registro de receitas e despesas, convites para grupos e geração de insights sobre os gastos.

## Funcionalidades Principais

- **Autenticação**: login por credenciais ou Google e geração/renovação de tokens.
- **Usuários**: consulta do usuário logado.
- **Grupos**: criação de grupos, envio de convites e consulta de saldo do grupo.
- **Categorias**: cadastro de categorias e grupos de categorias de despesas.
- **Despesas**: criação, atualização, exclusão e agendamento de despesas, além de geração de relatórios por período e categoria.
- **Receitas**: registro de receitas e consulta mensal.

## Arquitetura

O projeto utiliza o micro framework **Eon** (pasta `internal/pkg/eon`) para gerenciar ciclos de vida e injeção de dependências. Cada domínio possui um módulo próprio sob `internal/modules`, contendo:

- `controller/` – handlers HTTP definidos com o framework Fiber.
- `usecase/` – regras de negócio.
- `postgres/` – repositórios com implementação para Postgres.
- `module/` – arquivo que registra o módulo no Eon para bootstrapping.

Módulos comuns (servidor HTTP, configuração, banco de dados) ficam em `internal/pkg`. Há ainda um módulo `internal/shared` que provê clientes de infraestrutura (JWT, envio de e-mail, pub/sub) e middlewares.

As migrações do banco estão em `database/migrations` e podem ser executadas pelo script `database/migrate.sh`.

## Execução Local

1. Inicie o Postgres via Docker Compose:
   ```bash
   docker compose up db -d
   ```
2. Aplique as migrações:
   ```bash
   ./database/migrate.sh up ./database/migrations
   ```
3. Inicie a aplicação:
   ```bash
   ENV=development go run main.go
   ```

Há também um `Makefile` com atalhos para essas tarefas.

## Testes e Qualidade de Código

O repositório possui workflows do GitHub Actions para lint e testes (diretório `.github/workflows`). Para rodar manualmente:

```bash
go test ./...
```

## Deploy

O workflow `push-branch-main.yml` constrói a imagem Docker e faz o deploy no Google Cloud Run, executando previamente as migrações do banco de dados.


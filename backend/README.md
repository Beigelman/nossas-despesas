# Backend API - Nossas Despesas

The backend API is a RESTful service built with Go, providing the core functionality for expense management, group collaboration, and user authentication.

## Architecture

### Framework: Eon

The backend uses **Eon**, a custom micro framework located in `internal/pkg/eon`. Eon provides:

- **Dependency Injection Container**: Manages service dependencies and their lifecycle
- **Lifecycle Management**: Handles application bootstrapping, startup, and graceful shutdown
- **Module System**: Encapsulates domain logic into independent, composable modules

### Module Structure

Each domain module follows a consistent structure under `internal/modules/`:

```
module-name/
├── controller/     # HTTP handlers (Fiber routes)
├── usecase/       # Business logic layer
├── postgres/      # Data access layer (repositories)
└── module/        # Module registration for Eon
```

### Domain Modules

- **auth**: User authentication (credentials and Google OAuth), JWT token management
- **user**: User profile management
- **group**: Group creation, invitations, and balance calculations
- **category**: Category and category group management
- **expense**: Expense CRUD operations, scheduled expenses, reports, and split calculations
- **income**: Income registration and monthly queries

### Shared Infrastructure

Located in `internal/pkg/` and `internal/shared/`:

- **api**: HTTP server setup with Fiber framework
- **config**: Configuration management with environment variables
- **db**: Database connection and migration utilities
- **di**: Dependency injection container
- **eon**: Application framework
- **jwt**: JWT token generation and validation
- **email**: Email sending via Resend API
- **pubsub**: Pub/Sub messaging (Watermill)
- **logger**: Structured logging with slog
- **validator**: Request validation
- **predict**: ML service client for category prediction

## Technology Stack

- **Language**: Go 1.25+
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL with pgx/v5 driver
- **Migrations**: golang-migrate
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Validation**: go-playground/validator/v10
- **Testing**: testify, testcontainers-go
- **Logging**: slog (structured logging)
- **Error Tracking**: Sentry
- **Email**: Resend API
- **Pub/Sub**: Watermill

## Project Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── modules/             # Domain modules
│   │   ├── auth/
│   │   ├── user/
│   │   ├── group/
│   │   ├── category/
│   │   ├── expense/
│   │   └── income/
│   ├── pkg/                 # Shared packages
│   │   ├── api/             # HTTP server
│   │   ├── config/          # Configuration
│   │   ├── db/              # Database utilities
│   │   ├── di/              # Dependency injection
│   │   ├── eon/             # Application framework
│   │   ├── jwt/             # JWT utilities
│   │   ├── email/           # Email service
│   │   ├── pubsub/          # Pub/Sub messaging
│   │   └── ...
│   └── shared/              # Shared infrastructure
│       ├── middleware/      # HTTP middlewares
│       └── service/         # Shared services
├── database/
│   ├── migrations/         # SQL migration files
│   ├── schema.hcl          # Atlas schema definition
│   └── migrate.sh          # Migration script
├── scripts/                 # Utility scripts
│   ├── createusers/        # User creation script
│   ├── importincomes/      # Income import script
│   └── importsplit/        # SplitWise import script
├── templates/              # Email templates
├── config.go               # Configuration struct
├── go.mod                  # Go dependencies
└── Makefile                # Build and development commands
```

## Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 12+
- Docker and Docker Compose (for local database)

### Environment Variables

Create a `.env` file in the `backend/` directory:

```env
# Service Configuration
SERVICE_NAME=nossas-despesas-api
PORT=8080
ENV=development
LOG_LEVEL=info

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=app
DB_USER=root
DB_PASSWORD=root
# Or use connection string:
# DB_CONNECTION_STRING=postgres://user:password@localhost:5432/dbname?sslmode=disable

# JWT
JWT_SECRET=your-secret-key-here

# Email (Resend)
MAIL_API_KEY=your-resend-api-key
MAIL_SANDBOX_ID=your-sandbox-id

# ML Service
PREDICT_URL=http://localhost:8000

# Error Tracking (optional)
SENTRY_DSN=your-sentry-dsn
```

### Running Locally

1. **Start PostgreSQL**:
   ```bash
   make db
   # Or manually:
   docker compose up db -d
   ```

2. **Run Migrations**:
   ```bash
   make migrate-up
   # Or manually:
   ./database/migrate.sh up ./database/migrations
   ```

3. **Start the Server**:
   ```bash
   make dev
   # Or manually:
   ENV=development go run cmd/main.go
   ```

The API will be available at `http://localhost:8080`.

## Development

### Makefile Commands

```bash
# Database
make db                    # Start PostgreSQL container
make migrate-up           # Apply migrations
make migrate-down         # Rollback migrations
make migrate-new NAME=xxx # Create new migration
make migrate-diff NAME=xxx # Generate migration diff

# Development
make dev                  # Start development server
make format               # Format code
make lint                 # Run linter
make mock                 # Generate mocks

# Testing
make unit                 # Run unit tests
make integration          # Run integration tests
make test                 # Run all tests

# Scripts
make create-users         # Create test users
make import-incomes       # Import income data
make import-split         # Import from SplitWise
make reset-app            # Reset app (db + users + data)
```

### Code Organization

#### Controllers

Controllers handle HTTP requests and responses. They should:
- Parse and validate request data
- Call use cases
- Return appropriate HTTP responses
- Handle errors and map them to HTTP status codes

Example:
```go
func (c *Controller) CreateExpense(ctx *fiber.Ctx) error {
    var req CreateExpenseRequest
    if err := ctx.BodyParser(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "invalid request")
    }
    
    expense, err := c.createExpenseUseCase.Execute(ctx.Context(), req)
    if err != nil {
        return err
    }
    
    return ctx.Status(fiber.StatusCreated).JSON(expense)
}
```

#### Use Cases

Use cases contain business logic. They should:
- Be independent of HTTP concerns
- Validate business rules
- Coordinate between repositories
- Return domain errors

#### Repositories

Repositories handle data persistence. They should:
- Abstract database operations
- Map between domain models and database models
- Handle transactions when needed

### Database Migrations

Migrations use `golang-migrate` and are located in `database/migrations/`.

#### Creating a Migration

```bash
make migrate-new NAME=add_user_table
```

This creates two files:
- `YYYYMMDDHHMMSS_add_user_table.up.sql` - Migration up
- `YYYYMMDDHHMMSS_add_user_table.down.sql` - Migration down

#### Running Migrations

```bash
make migrate-up    # Apply all pending migrations
make migrate-down   # Rollback last migration
```

### Testing

#### Unit Tests

Test individual functions and use cases:
```bash
make unit
```

#### Integration Tests

Test database operations with testcontainers:
```bash
make integration
```

#### Running Specific Tests

```bash
go test -v ./internal/modules/expense/usecase/...
go test -v ./internal/modules/expense/postgres/...
```

## API Endpoints

### Authentication
- `POST /auth/sign-up` - Register with credentials
- `POST /auth/sign-in` - Login with credentials
- `POST /auth/sign-in/google` - Login with Google OAuth
- `POST /auth/refresh-token` - Refresh JWT token

### Users
- `GET /users/me` - Get current user

### Groups
- `POST /groups` - Create group
- `GET /groups/:id` - Get group details
- `POST /groups/:id/invite` - Invite user to group
- `POST /groups/:id/invite/accept` - Accept group invitation
- `GET /groups/:id/balance` - Get group balance

### Categories
- `GET /categories` - List all categories
- `POST /categories` - Create category
- `POST /categories/groups` - Create category group

### Expenses
- `GET /expenses` - List expenses
- `GET /expenses/:id` - Get expense details
- `POST /expenses` - Create expense
- `PUT /expenses/:id` - Update expense
- `DELETE /expenses/:id` - Delete expense
- `POST /expenses/scheduled` - Create scheduled expense
- `POST /expenses/scheduled/:id/generate` - Generate expenses from scheduled
- `GET /expenses/reports/period` - Get expenses by period
- `GET /expenses/reports/category` - Get expenses by category
- `POST /expenses/:id/recalculate-split` - Recalculate expense split
- `POST /expenses/predict-category` - Predict expense category (ML)

### Income
- `GET /income` - List income
- `POST /income` - Create income entry
- `GET /income/monthly` - Get monthly income

## Deployment

The backend is deployed to Google Cloud Run. The deployment workflow is triggered automatically on pushes to the main branch when backend files change.

### Environment Variables for Production

Ensure all required environment variables are set in Cloud Run:
- Database connection string
- JWT secret
- Email API keys
- ML service URL
- Sentry DSN (optional)

## Contributing

1. Follow Go best practices and the project's code style
2. Write tests for new features
3. Update documentation as needed
4. Run `make lint` before committing
5. Ensure all tests pass with `make test`


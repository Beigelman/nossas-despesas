# Nossas Despesas - Monorepo

This monorepo contains the services for the Nossas Despesas system, designed to organize and split expenses for couples:

- **Backend API** (Go): Main API with authentication, groups, expenses, and income management
- **ML Service** (Python): Machine learning service for automatic expense category classification
- **Web App** (Next.js): Frontend web application built with Next.js 15, React 19, and TypeScript

## Main Features

- **Authentication**: Login with credentials or Google OAuth, token generation and renewal
- **Users**: Query logged-in user information
- **Groups**: Create groups, send invites, and query group balance
- **Categories**: Register categories and category groups for expenses
- **Expenses**: Create, update, delete, and schedule expenses, plus generate reports by period and category
- **Income**: Register income and monthly queries

## Monorepo Structure

```
/
├── backend/              # Backend API (Go)
│   ├── internal/         # Internal backend code
│   ├── database/         # Database migrations and scripts
│   ├── templates/        # Email templates
│   └── scripts/          # Utility scripts
├── machine_learn/        # ML Service (Python)
│   ├── src/              # ML API source code
│   ├── training/         # Training scripts
│   └── models/           # Trained models
├── web/                  # Web App (Next.js)
│   ├── src/              # Application source code
│   ├── public/           # Static files
│   └── vercel.json       # Vercel configuration
└── .github/workflows/    # CI/CD separated by service
```

## Module Documentation

For detailed information about each module, please refer to their respective README files:

- **[Backend API](./backend/README.md)**: Complete guide to the Go backend, including architecture, API endpoints, development setup, and deployment
- **[Web App](./web/README.md)**: Frontend application documentation with Next.js setup, component structure, and features
- **[ML Service](./machine_learn/README.md)**: Machine learning service documentation covering model training, API endpoints, and prediction pipeline

# Web App - Nossas Despesas

The web application is a modern, responsive Progressive Web App (PWA) built with Next.js 15, React 19, and TypeScript, providing an intuitive interface for managing expenses, groups, and financial tracking.

## Technology Stack

- **Framework**: Next.js 15 (App Router)
- **React**: React 19 (RC)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: Radix UI
- **State Management**: 
  - Zustand (client state)
  - TanStack Query (server state)
- **Forms**: React Hook Form + Zod validation
- **Authentication**: NextAuth v4 (with Google OAuth and credentials)
- **HTTP Client**: Axios
- **Charts**: Recharts
- **Date Handling**: date-fns
- **PWA**: @ducanh2912/next-pwa
- **Package Manager**: pnpm

## Project Structure

```
web/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── (application)/      # Protected routes (require auth)
│   │   │   ├── expenses/       # Expense management pages
│   │   │   ├── groups/         # Group management pages
│   │   │   ├── income/         # Income pages
│   │   │   └── ...
│   │   ├── (auth)/             # Authentication routes
│   │   │   ├── signin/         # Sign in page
│   │   │   └── ...
│   │   ├── api/                # API routes
│   │   ├── layout.tsx          # Root layout
│   │   └── page.tsx            # Home page
│   ├── components/             # React components
│   │   ├── ui/                 # Reusable UI components (Radix UI)
│   │   ├── app-header.tsx      # Application header
│   │   ├── side-bar.tsx        # Sidebar navigation
│   │   ├── search-bar.tsx      # Search functionality
│   │   └── ...
│   ├── domain/                 # TypeScript domain types
│   │   ├── expense.d.ts
│   │   ├── category.ts
│   │   ├── group.d.ts
│   │   └── user.d.ts
│   ├── hooks/                  # Custom React hooks
│   │   ├── use-expenses.ts     # Expense data fetching
│   │   ├── use-balance.ts      # Balance calculations
│   │   ├── use-categories.ts   # Category management
│   │   └── ...
│   ├── http/                   # HTTP client configuration
│   │   ├── api.ts              # Axios instance
│   │   └── interceptors.ts     # Request/response interceptors
│   ├── lib/                    # Utility functions
│   │   ├── auth/               # Authentication utilities
│   │   ├── utils.ts            # General utilities
│   │   └── ...
│   ├── providers/              # React context providers
│   │   ├── query-provider.tsx  # TanStack Query provider
│   │   ├── theme-provider.tsx  # Theme (dark/light mode)
│   │   └── ...
│   ├── schemas/                # Zod validation schemas
│   ├── styles/                 # Global styles
│   │   └── globals.css         # Tailwind CSS imports
│   └── types/                  # TypeScript type definitions
├── public/                     # Static assets
│   ├── icons/                # PWA icons
│   ├── manifest.json          # PWA manifest
│   └── sw.js                  # Service worker
├── .next/                      # Next.js build output (gitignored)
├── next.config.js             # Next.js configuration
├── tailwind.config.js          # Tailwind CSS configuration
├── tsconfig.json               # TypeScript configuration
└── package.json                # Dependencies and scripts
```

## Getting Started

### Prerequisites

- Node.js 22.x
- pnpm (recommended) or npm/yarn

### Installation

1. **Install dependencies**:
   ```bash
   cd web
   pnpm install
   ```

2. **Configure environment variables**:
   
   Create a `.env.local` file in the `web/` directory:
   ```env
   # Public variables (exposed to browser)
   NEXT_PUBLIC_API_URL=http://localhost:8080
   NEXT_PUBLIC_BASE_URL=http://localhost:3000

   # Authentication
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   NEXTAUTH_SECRET=your-nextauth-secret
   NEXTAUTH_URL=http://localhost:3000
   ```

   To generate `NEXTAUTH_SECRET`:
   ```bash
   openssl rand -base64 32
   ```

3. **Start the development server**:
   ```bash
   pnpm dev
   ```

   The application will be available at `http://localhost:3000`.

## Development

### Available Scripts

```bash
pnpm dev          # Start development server
pnpm build        # Build for production
pnpm start        # Start production server
pnpm lint         # Run ESLint and fix issues
```

### Code Organization

#### Pages (App Router)

Pages are located in `src/app/` using Next.js App Router:
- `(application)/` - Protected routes that require authentication
- `(auth)/` - Public authentication pages
- `api/` - API routes (if needed)

#### Components

- **UI Components** (`components/ui/`): Reusable, accessible components built on Radix UI
- **Feature Components** (`components/`): Domain-specific components like `app-header`, `side-bar`, etc.

#### Hooks

Custom hooks in `hooks/` encapsulate data fetching and business logic:
- `use-expenses.ts` - Fetch and manage expenses
- `use-balance.ts` - Calculate user/group balances
- `use-categories.ts` - Manage categories
- `use-group.ts` - Group operations
- `use-income.ts` - Income management
- `use-predict.ts` - ML category prediction

#### HTTP Client

The HTTP client is configured in `src/http/`:
- `api.ts` - Axios instance with base configuration
- `interceptors.ts` - Request/response interceptors for auth tokens, error handling

#### Domain Types

TypeScript domain types are defined in `src/domain/` to ensure type safety across the application.

### Styling

The project uses **Tailwind CSS** for styling with:
- Custom theme configuration in `tailwind.config.js`
- Dark mode support via `next-themes`
- Responsive design utilities
- Custom component classes

### Authentication

Authentication is handled by **NextAuth v4** with:
- **Credentials Provider**: Email/password authentication
- **Google Provider**: OAuth 2.0 with Google
- Session management with JWT tokens
- Protected routes via middleware

### State Management

- **TanStack Query**: Server state (API data, caching, refetching)
- **Zustand**: Client state (UI state, preferences)
- **React Context**: Theme provider, query provider

### Progressive Web App (PWA)

The app is configured as a PWA with:
- Service worker for offline support
- App manifest for installability
- Icons for various platforms
- Caching strategies for assets

## Features

### Expense Management
- Create, edit, and delete expenses
- Schedule recurring expenses
- Split expenses among group members
- Filter and search expenses
- View expenses by category or period
- Automatic category prediction (ML integration)

### Group Management
- Create and manage groups
- Invite users to groups
- View group balance and member contributions
- Track shared expenses

### Income Tracking
- Register income entries
- View monthly income summaries
- Track income over time

### Categories
- Manage expense categories
- Organize categories into groups
- Custom category icons

### Reports and Insights
- Expense reports by period
- Category-based expense analysis
- Balance calculations
- Visual charts and graphs

## Building for Production

```bash
pnpm build
```

This creates an optimized production build in `.next/`.

## Deployment

The web app is deployed to **Vercel**. The deployment workflow:
- Automatically builds on push to main branch
- Runs linting and type checking
- Deploys to Vercel preview/production

### Vercel Configuration

Ensure all environment variables are set in Vercel:
- Go to project settings → Environment Variables
- Add all variables from `.env.local`
- Configure for Production, Preview, and Development environments


## Code Quality

### Linting

```bash
pnpm lint
```

Uses ESLint with Next.js and custom configurations.

### Type Checking

```bash
pnpm exec tsc --noEmit
```

## Contributing

1. Follow the existing code style and patterns
2. Use TypeScript for all new code
3. Write reusable components when possible
4. Update types in `domain/` when adding new features
5. Run `pnpm lint` before committing
6. Ensure the app builds successfully with `pnpm build`

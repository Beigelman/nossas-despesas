# Guia de Migração: NextAuth v4 → v5 (Auth.js)

## 📋 Resumo da Migração

A migração do NextAuth v4 para v5 (Auth.js) é uma **reescrita significativa** que resolve os problemas de compatibilidade com Next.js 15, mas requer mudanças substanciais no código.

## 🔄 Principais Mudanças

### 1. Estrutura de Arquivos

**v4:**
```
src/app/api/auth/[...nextauth]/route.ts
src/lib/auth/next-auth-options.ts
```

**v5:**
```
src/auth.ts (ou src/auth.config.ts)
src/app/api/auth/[...route]/route.ts
```

### 2. Configuração (auth.ts)

**v4:**
```typescript
import { AuthOptions } from 'next-auth'
export const nextAuthOptions: AuthOptions = { ... }
```

**v5:**
```typescript
import NextAuth from "next-auth"
import Google from "next-auth/providers/google"
import Credentials from "next-auth/providers/credentials"

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    Google({ ... }),
    Credentials({ ... })
  ],
  callbacks: { ... }
})
```

### 3. API Route Handler

**v4:**
```typescript
import NextAuth from 'next-auth'
import { nextAuthOptions } from '@/lib/auth/next-auth-options'

const handler = NextAuth(nextAuthOptions)
export { handler as GET, handler as POST }
```

**v5:**
```typescript
import { handlers } from "@/auth"
export const { GET, POST } = handlers
```

### 4. Server Components (getServerSession)

**v4:**
```typescript
import { getServerSession } from 'next-auth'
import { nextAuthOptions } from '@/lib/auth/next-auth-options'

const session = await getServerSession(nextAuthOptions)
```

**v5:**
```typescript
import { auth } from "@/auth"

const session = await auth()
```

### 5. Client Components (useSession)

**v4:**
```typescript
import { useSession, signIn, signOut } from 'next-auth/react'

const { data: session } = useSession()
await signIn('google')
await signOut()
```

**v5:**
```typescript
import { useSession } from "next-auth/react"
import { signIn, signOut } from "@/auth"

const { data: session } = useSession()
await signIn('google')
await signOut()
```

### 6. SessionProvider

**v4:**
```typescript
import { SessionProvider } from 'next-auth/react'
<SessionProvider>{children}</SessionProvider>
```

**v5:**
```typescript
import { SessionProvider } from "next-auth/react"
// Mantém a mesma API, mas agora usa o auth() do v5 internamente
<SessionProvider>{children}</SessionProvider>
```

### 7. Callbacks e Types

**v4:**
```typescript
callbacks: {
  async jwt({ token, user, account }) { ... },
  async session({ session, token }) { ... }
}
```

**v5:**
```typescript
callbacks: {
  async jwt({ token, user, account }) { ... },
  async session({ session, token }) { ... }
}
// Similar, mas com tipos mais rigorosos
```

## 📝 Arquivos que Precisam ser Modificados

### Arquivos Principais:
1. ✅ `src/lib/auth/next-auth-options.ts` → `src/auth.ts` (reescrito completo)
2. ✅ `src/app/api/auth/[...nextauth]/route.ts` → `src/app/api/auth/[...route]/route.ts`
3. ✅ `src/lib/auth/get-server-session.ts` → Removido (usar `auth()` diretamente)
4. ✅ `src/providers/session.tsx` → Pode manter similar
5. ✅ `src/hooks/use-user.ts` → Atualizar imports
6. ✅ `src/http/private-client.ts` → Atualizar imports
7. ✅ Todos os layouts que usam `getServerSession` → Usar `auth()` do v5

### Arquivos com Mudanças Menores:
- `src/components/signin-with-google-button.tsx`
- `src/components/user-nav.tsx`
- `src/components/signout-button.tsx`
- `src/app/(application)/group/[inviteId]/accept/refreshing-session.tsx`
- `src/app/(auth)/register/page.tsx`

## 🚀 Passos da Migração

### Passo 1: Instalar NextAuth v5
```bash
npm install next-auth@beta
# ou
npm install next-auth@5.0.0-beta.29
```

### Passo 2: Criar novo arquivo de configuração (src/auth.ts)
- Migrar providers (Google, Credentials)
- Migrar callbacks (jwt, session)
- Ajustar tipos e interfaces

### Passo 3: Atualizar API Route
- Renomear `[...nextauth]` para `[...route]`
- Usar handlers exportados do auth.ts

### Passo 4: Atualizar Server Components
- Substituir `getServerSession(nextAuthOptions)` por `auth()`
- Remover helper `get-server-session.ts`

### Passo 5: Atualizar Client Components
- Manter `useSession` do `next-auth/react`
- Atualizar imports de `signIn`/`signOut` se necessário

### Passo 6: Atualizar Types
- Criar tipos customizados para session se necessário
- Ajustar tipos de User, Token, etc.

### Passo 7: Testar
- Testar login com Google
- Testar refresh token
- Testar todas as rotas protegidas
- Testar logout

## ⚠️ Pontos de Atenção

1. **Credentials Provider**: Pode precisar de ajustes na forma como funciona
2. **Custom Session Types**: Pode precisar ajustar tipos customizados
3. **Refresh Token Logic**: Verificar se a lógica de refresh ainda funciona
4. **Middleware**: Se houver middleware customizado, pode precisar ajustes
5. **Environment Variables**: Verificar se todas as variáveis estão corretas

## 📚 Recursos

- [Documentação NextAuth v5](https://authjs.dev/)
- [Guia de Migração Oficial](https://authjs.dev/getting-started/migrating-to-v5)
- [Exemplos v5](https://authjs.dev/getting-started/example)

## 💡 Recomendação

**Opção 1: Migrar Agora (Recomendado se você tem tempo)**
- Resolve os erros de compatibilidade com Next.js 15
- API mais moderna e melhor suportada
- Requer tempo para testar tudo

**Opção 2: Aguardar (Se você precisa de estabilidade)**
- Os erros são apenas avisos em desenvolvimento
- A aplicação funciona normalmente
- Pode aguardar versão estável do v5

## 🎯 Estimativa de Tempo

- **Migração básica**: 2-4 horas
- **Testes e ajustes**: 2-4 horas
- **Total**: 4-8 horas de trabalho


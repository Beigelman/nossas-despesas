import { User } from '@/domain/user'

declare module 'next-auth' {
  interface Session {
    user: User
    token: string
    refreshToken: string
  }
}

import { AuthOptions } from 'next-auth'
import CredentialsProvider from 'next-auth/providers/credentials'
import GoogleProvider from 'next-auth/providers/google'
import { toast } from 'sonner'

import { LoginUser, refreshToken, signInWithGoogle } from '@/http/auth'

export const nextAuthOptions: AuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID ?? '',
      clientSecret: process.env.GOOGLE_CLIENT_SECRET ?? '',
    }),
    CredentialsProvider({
      id: 'refresh-token',
      credentials: {
        refreshToken: {},
      },
      async authorize(credentials) {
        const resp = await refreshToken(credentials?.refreshToken ?? '')
        if (!resp) {
          return null
        }

        return resp as never
      },
    }),
  ],
  callbacks: {
    async session({ session, token }) {
      const userInfo = token.session as LoginUser
      session.user = userInfo.user
      session.token = userInfo.token
      session.refreshToken = userInfo.refreshToken
      return session
    },
    async jwt({ token, user, account }) {
      if (account?.provider === 'google' && !token.session) {
        try {
          const resp = await signInWithGoogle(account.id_token ?? '')
          token.session = resp
        } catch (error) {
          toast.error('Erro ao logar com Google')
        }

        return token
      }

      if (user && !token.session) {
        token.session = user
        return token
      }

      return token
    },
  },
}

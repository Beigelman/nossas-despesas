import NextAuth from 'next-auth'

import { nextAuthOptions } from '@/lib/auth/next-auth-options'

const handler = NextAuth(nextAuthOptions)

export async function GET(request: Request, context: { params: Promise<{ nextauth: string[] }> }) {
  const resolvedParams = await context.params
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  return handler(request, { params: resolvedParams } as any)
}

export async function POST(request: Request, context: { params: Promise<{ nextauth: string[] }> }) {
  const resolvedParams = await context.params
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  return handler(request, { params: resolvedParams } as any)
}

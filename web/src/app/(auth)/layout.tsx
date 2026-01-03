import { headers } from 'next/headers'
import { redirect } from 'next/navigation'

import { getServerSession } from '@/lib/auth/get-server-session'

type AuthLayoutProps = {
  children: React.ReactNode
}

type HeaderList = Awaited<ReturnType<typeof headers>>

export default async function AuthLayout({ children }: AuthLayoutProps) {
  const session = await getServerSession()
  const headersList = await headers()
  const previousPath = getPreviousPath(headersList)

  if (session) {
    if (previousPath) {
      redirect(previousPath)
    } else {
      redirect('/expenses')
    }
  }

  return <>{children}</>
}

function getPreviousPath(headers: HeaderList): string {
  try {
    const urlObject = new URL(headers.get('x-url') || '')
    const previous = urlObject.searchParams.get('previous') || ''
    const previousURL = new URL(previous)
    return previousURL.pathname
  } catch {
    return ''
  }
}

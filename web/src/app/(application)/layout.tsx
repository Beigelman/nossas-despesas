import { headers } from 'next/headers'
import { redirect } from 'next/navigation'

import { AppHeader } from '@/components/app-header'
import { getServerSession } from '@/lib/auth/get-server-session'

type PrivateLayoutProps = {
  children: React.ReactNode
}

export default async function PrivateLayout({ children }: PrivateLayoutProps) {
  const session = await getServerSession()
  const headersList = await headers()
  const headerURL = headersList.get('x-url') || ''

  if (!session) {
    redirect('/login?previous=' + headerURL)
  }

  return (
    <div vaul-drawer-wrapper="">
      <div className="relative flex min-h-screen flex-col bg-background">
        <AppHeader />
        <main className="flex flex-1">{children}</main>
      </div>
    </div>
  )
}

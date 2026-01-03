import '@/styles/globals.css'

import { Metadata, Viewport } from 'next'
import { headers } from 'next/headers'
import { redirect } from 'next/navigation'

import { Toaster } from '@/components/ui/sonner'
import { getServerSession } from '@/lib/auth/get-server-session'
import { Providers } from '@/providers'

type RootLayoutProps = {
  children: React.ReactNode
}

export const metadata: Metadata = {
  title: 'Nossas despesas',
  description: 'Seu app the controle de despesas com o seu mozão',
  generator: 'Next.js',
  manifest: '/manifest.json',
  keywords: [
    'controle de despesas',
    'divisor de contas',
    'despesas casal',
    'finanças',
    'dinheiro casal',
    'nossas despesas',
  ],
  authors: [
    {
      name: 'Daniel Beigelman',
      url: 'https://www.linkedin.com/in/daniel-beigelman/',
    },
  ],
  icons: [
    { rel: 'apple-touch-icon', url: 'icons/apple-touch-icon.png' },
    { rel: 'icon', url: '/icons/icon-rounded.png' },
  ],
}

export const viewport: Viewport = {
  initialScale: 1,
  maximumScale: 1,
  minimumScale: 1,
  width: 'device-width',
}

export default async function RootLayout({ children }: RootLayoutProps) {
  const session = await getServerSession()
  const headersList = await headers()
  const headerURL = headersList.get('x-url') || ''

  let pathname = '/'
  try {
    pathname = new URL(headerURL).pathname
  } catch {
    // noop - keep default pathname
  }

  if (session && pathname === '/') {
    redirect('/expenses')
  }

  return (
    <html lang="en" suppressHydrationWarning>
      <head />
      <body className="h-screen w-full">
        <Providers>{children}</Providers>
        <Toaster />
      </body>
    </html>
  )
}

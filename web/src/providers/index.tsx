import { TooltipProvider } from '@radix-ui/react-tooltip'
import { NextIntlClientProvider } from 'next-intl'
import { getLocale, getMessages } from 'next-intl/server'

import QueryProvider from './query'
import NextAuthSessionProvider from './session'
import { ThemeProvider } from './theme-provider'

type ProvidersProps = {
  children: React.ReactNode
}

async function Providers({ children }: ProvidersProps): Promise<React.ReactNode> {
  const locale = await getLocale()
  const messages = await getMessages()

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <TooltipProvider>
        <QueryProvider>
          <NextAuthSessionProvider>
            <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
              {children}
            </ThemeProvider>
          </NextAuthSessionProvider>
        </QueryProvider>
      </TooltipProvider>
    </NextIntlClientProvider>
  )
}

export { Providers }

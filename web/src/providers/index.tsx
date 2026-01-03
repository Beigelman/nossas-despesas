import { TooltipProvider } from '@radix-ui/react-tooltip'

import QueryProvider from './query'
import NextAuthSessionProvider from './session'
import { ThemeProvider } from './theme-provider'

type ProvidersProps = {
  children: React.ReactNode
}
function Providers({ children }: ProvidersProps): React.ReactNode {
  return (
    <TooltipProvider>
      <QueryProvider>
        <NextAuthSessionProvider>
          <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
            {children}
          </ThemeProvider>
        </NextAuthSessionProvider>
      </QueryProvider>
    </TooltipProvider>
  )
}

export { Providers }

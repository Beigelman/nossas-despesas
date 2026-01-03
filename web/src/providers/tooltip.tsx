import { TooltipProvider } from '@/components/ui/tooltip'

type TooltipProviderProps = {
  children: React.ReactNode
}

export default function GlobalTooltipProvider({ children }: TooltipProviderProps) {
  return <TooltipProvider>{children}</TooltipProvider>
}

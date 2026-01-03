import { cn } from '@/lib/utils'

type SpinningProps = React.HTMLAttributes<HTMLDivElement>

function Spinning({ className, ...props }: SpinningProps) {
  return (
    <div
      className={cn('inline-block h-10 w-10 animate-spin rounded-full border-4 border-t-transparent', className)}
      {...props}
    />
  )
}

export { Spinning }

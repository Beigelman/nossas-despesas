import { Wand2 } from 'lucide-react'

import { cn } from '@/lib/utils'

import { Icon } from './ui/icon'

type CategoryIconProps = {
  name: string
  size?: number
  className?: string
  hover?: boolean
  isPredicted?: boolean
  isPredicting?: boolean
  blendBackground?: boolean
}

function CategoryIcon({ name, size, hover, className, isPredicted, isPredicting, blendBackground }: CategoryIconProps) {
  const shouldGlow = isPredicting && !isPredicted

  return (
    <div
      className={cn(
        'relative z-10 p-[2px]', // espaço para a borda “vazar”
        shouldGlow && 'animate-glow-border',
      )}
    >
      <div
        className={cn(
          'relative z-10 flex items-center justify-center rounded-md',
          blendBackground ? 'bg-transparent' : 'bg-background',
          className,
          hover && 'hover:bg-accent hover:text-accent-foreground',
        )}
      >
        <Icon name={name} size={size ?? 18} />
      </div>

      {isPredicted && !isPredicting && (
        <span className="pointer-events-none absolute -bottom-1 -right-1 z-20 flex h-5 w-5 items-center justify-center rounded-full bg-sky-500 text-white shadow-md">
          <Wand2 className="h-3 w-3" />
        </span>
      )}
    </div>
  )
}

export { CategoryIcon }

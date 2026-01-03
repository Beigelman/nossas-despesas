import { X as XIcon } from 'lucide-react'

import { cn } from '@/lib/utils'

type XProps = React.ButtonHTMLAttributes<HTMLButtonElement>

function X({ className, ...props }: XProps) {
  return (
    <button className={cn('text-red-700 hover:border-b-2 hover:border-red-700', className)} {...props}>
      <XIcon size={18} />
    </button>
  )
}

export { X }

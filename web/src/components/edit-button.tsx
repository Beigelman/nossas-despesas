import { Edit as EditIcon } from 'lucide-react'

import { cn } from '@/lib/utils'

type EditProps = React.ButtonHTMLAttributes<HTMLButtonElement>

function Edit({ className, ...props }: EditProps) {
  return (
    <button className={cn('hover:border-b-2 hover:border-black dark:border-white', className)} {...props}>
      <EditIcon size={18} />
    </button>
  )
}

export { Edit }

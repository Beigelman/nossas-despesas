import { cn, formatCurrency } from '@/lib/utils'

import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar'

type UserBalanceProps = React.ButtonHTMLAttributes<HTMLDivElement> & {
  user: {
    name: string
    avatar: string
    balance: {
      amount: number
      type: 'income' | 'outcome' | 'neutral'
    }
  }
}

function getCurrencyColor(type: 'income' | 'outcome' | 'neutral'): string {
  switch (type) {
    case 'income':
      return 'text-green-900 dark:text-green-500'
    case 'outcome':
      return 'text-red-900 dark:text-red-500'
    case 'neutral':
      return 'text-gray-900 dark:text-gray-500'
  }
}

function UserBalance({ user, className }: UserBalanceProps) {
  const balance = formatCurrency(user.balance.amount)
  const fallback = user.name.substring(0, 2).toUpperCase()

  return (
    <div className={cn('flex', className)}>
      <Avatar>
        <AvatarImage src={user.avatar} alt={user.name} />
        <AvatarFallback>{fallback}</AvatarFallback>
      </Avatar>
      <div className="ml-2 flex flex-col">
        <span>{user.name}</span>
        <p className={cn(getCurrencyColor(user.balance.type), '-mt-1 text-sm')}>
          {`${user.balance.type === 'income' ? 'receberá' : 'deve'} ${balance}`}
        </p>
      </div>
    </div>
  )
}

export { UserBalance }

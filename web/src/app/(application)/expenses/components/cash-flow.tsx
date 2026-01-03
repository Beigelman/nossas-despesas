import { cn } from '@/lib/utils'

type CashFlowProps = {
  className?: string
  amount: string
  description: string
  cashFlowType: 'income' | 'outcome' | 'neutral'
}

function CashFlow(props: CashFlowProps) {
  let cashFlowColor: string
  switch (props.cashFlowType) {
    case 'income':
      cashFlowColor = 'text-green-900 dark:text-green-500'
      break
    case 'outcome':
      cashFlowColor = 'text-red-900 dark:text-red-500'
      break
    default:
      cashFlowColor = 'text-black dark:text-white'
      break
  }

  return (
    <div className={cn('flex gap-2 overflow-hidden whitespace-nowrap text-sm text-slate-400', props.className)}>
      {props.description}
      <h3 className={`font-bold ${cashFlowColor}`}>{props.amount}</h3>
    </div>
  )
}

export { CashFlow }

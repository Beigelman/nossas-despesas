import { DateDisplay } from '@/components/date-display'
import { Button } from '@/components/ui/button'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { Expense } from '@/domain/expense'
import { useCategories } from '@/hooks/use-categories'
import { useGroup } from '@/hooks/use-group'
import { formatCurrency } from '@/lib/utils'

import { CategoryIcon } from '../../../../components/category-icon'
import { Separator } from '../../../../components/ui/separator'
import { CashFlow } from './cash-flow'
import { DeleteExpenseButton } from './delete-expense-button'
import { ExpenseDetails } from './expense-details'
import { SaveExpenseButton } from './save-expense-button'
import { SplitRatio } from './split-ratio'

type ExpenseLineProps = React.ButtonHTMLAttributes<HTMLDivElement> & {
  expense: Expense
  selectedExpense: number
  setSelectedExpense: (open: number) => void
}

function ExpenseLine({ expense, setSelectedExpense, selectedExpense }: ExpenseLineProps) {
  const { me, partner } = useGroup()
  const { getCategory } = useCategories()

  const partnerName = partner?.name.split(' ')[0]
  const payerName = expense.payerId === me?.id ? 'você' : partnerName
  const realAmount = expense.amount - (expense.refundAmount ?? 0)
  const payerAmount = (realAmount * expense.splitRatio.payer) / 100
  const receiverAmount = realAmount - payerAmount

  return (
    <Collapsible
      onOpenChange={(open) => (open ? setSelectedExpense(expense.id) : setSelectedExpense(0))}
      open={selectedExpense === expense.id}
    >
      <CollapsibleTrigger className="relative flex w-full items-start gap-3 overflow-hidden p-4 transition-all hover:bg-gray-50 dark:hover:bg-gray-900">
        <div className="flex w-full justify-between gap-3">
          <div className="flex gap-2 overflow-hidden">
            <DateDisplay date={expense.createdAt} />
            <CategoryIcon name={getCategory(expense.categoryId)?.icon ?? ''} size={26} blendBackground />
            <span className="mt-1 truncate font-medium">{expense.name}</span>
          </div>
          <div className="flex flex-col items-end justify-center">
            <CashFlow amount={formatCurrency(realAmount)} description={`${payerName} pagou`} cashFlowType="neutral" />
            <CashFlow
              amount={formatCurrency(receiverAmount)}
              description="e emprestou"
              cashFlowType={expense.payerId === me?.id ? 'income' : 'outcome'}
            />
          </div>
        </div>
        <SplitRatio
          receiverRatio={expense.splitRatio.receiver}
          payerRatio={expense.splitRatio.payer}
          className="absolute bottom-2 left-14"
        />
      </CollapsibleTrigger>
      <CollapsibleContent>
        <div className="flex flex-col gap-3 px-4 pb-4 pt-1">
          <div className="flex items-center justify-end gap-2">
            <DeleteExpenseButton expense={expense}>
              <Button className="h-6" variant="destructive">
                Deletar
              </Button>
            </DeleteExpenseButton>
            <SaveExpenseButton expense={expense} type="update">
              <Button className="h-6">Editar</Button>
            </SaveExpenseButton>
          </div>
          <Separator />
          <ExpenseDetails expense={expense} />
        </div>
      </CollapsibleContent>
      <Separator />
    </Collapsible>
  )
}

export { ExpenseLine }

'use client'

import { EditIcon, SaveIcon } from 'lucide-react'
import { useState } from 'react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Expense } from '@/domain/expense'
import { useExpenses } from '@/hooks/use-expenses'
import { cn, moneyMask } from '@/lib/utils'

type RefundFormProps = React.HTMLAttributes<HTMLDivElement> & {
  expense: Expense
}

function RefundForm({ className, expense }: RefundFormProps) {
  const [refund, setRefund] = useState(expense.refundAmount?.toString() ?? '')
  const [isEditing, setIsEditing] = useState(false)
  const { updateExpense } = useExpenses('')

  function handleSave() {
    try {
      if (refund !== '') {
        updateExpense({ id: expense.id, refundAmount: parseInt(refund) })
      }
    } catch {
      setRefund('')
    } finally {
      setIsEditing(false)
    }
  }

  return (
    <div className={cn('flex justify-between gap-2', className)}>
      <div className="flex flex-1 flex-col justify-center">
        <span className="text-xs font-bold">Reembolso</span>
        <div className="flex flex-1 items-end gap-1">
          <span className="text-sm">R$</span>
          <Input
            className="p-0"
            disabled={!isEditing}
            value={moneyMask(refund)}
            placeholder="0,00"
            onChange={(e) => setRefund(e.target.value.replace(/\D/g, ''))}
            variant="outline"
            inputMode="decimal"
          />
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Button className="h-7 w-7 p-0" onClick={() => setIsEditing(true)} disabled={isEditing}>
          <EditIcon size={12} />
        </Button>
        <Button className="h-7 w-7 p-0" disabled={!isEditing} onClick={handleSave}>
          <SaveIcon size={12} />
        </Button>
      </div>
    </div>
  )
}

export { RefundForm }

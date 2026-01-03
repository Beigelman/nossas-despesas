import { useEffect, useState } from 'react'
import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { moneyMask } from '@/lib/utils'
import { expenseFormSchema } from '@/schemas'

type RefundFieldProps = {
  form: UseFormReturn<z.infer<typeof expenseFormSchema>>
}

function RefundField({ form }: RefundFieldProps) {
  const [hasRefund, setHasRefund] = useState(false)

  useEffect(() => {
    const state = form.getValues('refundAmount')
    setHasRefund(!!state)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  function handleRefund(checked: boolean) {
    if (!checked) {
      form.setValue('refundAmount', '')
    }
    setHasRefund(checked)
  }

  return (
    <div>
      <div className="flex items-center justify-center gap-2">
        <span className="text-sm">Tem reembolso?</span>
        <Switch checked={hasRefund} onCheckedChange={handleRefund} />
      </div>
      <FormField
        name="refundAmount"
        control={form.control}
        render={({ field }) => (
          <FormItem className={`pl-3 ${!hasRefund ? 'hidden' : ''}`}>
            <FormControl>
              <Input
                placeholder="R$ 0,00"
                inputMode="decimal"
                onChange={(event) => field.onChange(event.currentTarget.value.replace(/\D/g, ''))}
                value={moneyMask(field.value)}
                variant="outline"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    </div>
  )
}

export { RefundField }

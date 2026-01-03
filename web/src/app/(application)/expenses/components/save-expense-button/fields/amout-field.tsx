import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { moneyMask } from '@/lib/utils'
import { expenseFormSchema } from '@/schemas'

type AmountFieldProps = {
  form: UseFormReturn<z.infer<typeof expenseFormSchema>>
}

function AmountField({ form }: AmountFieldProps) {
  return (
    <FormField
      name="amount"
      control={form.control}
      render={({ field }) => (
        <FormItem className="flex flex-col">
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
  )
}

export { AmountField }

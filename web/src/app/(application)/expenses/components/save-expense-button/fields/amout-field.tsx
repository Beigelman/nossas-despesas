'use client'

import { useLocale } from 'next-intl'
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
  const locale = useLocale()
  const placeholder = locale === 'en' ? '$0.00' : 'R$ 0,00'

  return (
    <FormField
      name="amount"
      control={form.control}
      render={({ field }) => (
        <FormItem className="flex flex-col">
          <FormControl>
            <Input
              placeholder={placeholder}
              inputMode="decimal"
              onChange={(event) => field.onChange(event.currentTarget.value.replace(/\D/g, ''))}
              value={moneyMask(field.value, locale)}
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

import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { FormControl, FormField, FormItem } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { incomeFormSchema } from '@/schemas'

type AmountFieldProps = {
  form: UseFormReturn<z.infer<typeof incomeFormSchema>>
}
const moneyMask = (value = '') => {
  if (value === '') {
    return '0,00'
  }

  value = value.replace('.', '').replace(',', '').replace(/\D/g, '')
  const options = { minimumFractionDigits: 2 }
  return new Intl.NumberFormat('pt-BR', options).format(parseFloat(value) / 100)
}

function AmountField({ form }: AmountFieldProps) {
  return (
    <FormField
      name="amount"
      control={form.control}
      render={({ field }) => (
        <FormItem className="flex w-[240px] items-end">
          <span className="mb-1 mr-2">R$</span>
          <FormControl>
            <Input
              inputMode="decimal"
              placeholder="0.00"
              onChange={(event) => field.onChange(event.currentTarget.value.replace(/\D/g, ''))}
              value={moneyMask(field.value)}
              variant="outline"
            />
          </FormControl>
        </FormItem>
      )}
    />
  )
}

export { AmountField }

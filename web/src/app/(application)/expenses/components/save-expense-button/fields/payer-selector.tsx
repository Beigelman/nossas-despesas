import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { Badge } from '@/components/ui/badge'
import { FormControl, FormField, FormItem } from '@/components/ui/form'
import { CustomSelectTrigger, Select, SelectContent, SelectItem } from '@/components/ui/select'
import { User } from '@/domain/user'
import { expenseFormSchema } from '@/schemas'

type PayerSelectorProps = {
  form: UseFormReturn<z.infer<typeof expenseFormSchema>>
  user?: User
  partner?: User
  className?: string
}

function PayerSelector({ form, user, partner, className }: PayerSelectorProps) {
  return (
    <FormField
      name="payerId"
      control={form.control}
      render={({ field }) => (
        <FormItem>
          <Select onValueChange={(value) => field.onChange(parseInt(value))} defaultValue={String(field.value)}>
            <FormControl>
              <CustomSelectTrigger>
                <Badge className={className}>{String(field.value) === String(user?.id) ? 'você' : partner?.name}</Badge>
              </CustomSelectTrigger>
            </FormControl>
            <SelectContent>
              <SelectItem value={String(user?.id)}>você</SelectItem>
              <SelectItem value={String(partner?.id)}>{partner?.name}</SelectItem>
            </SelectContent>
          </Select>
        </FormItem>
      )}
    />
  )
}

export { PayerSelector }

import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { Badge } from '@/components/ui/badge'
import { FormControl, FormField, FormItem } from '@/components/ui/form'
import { CustomSelectTrigger, Select, SelectContent, SelectItem } from '@/components/ui/select'
import { expenseFormSchema } from '@/schemas'

type SplitRatioSelectorProps = {
  form: UseFormReturn<z.infer<typeof expenseFormSchema>>
}

function SplitRatioSelector({ form }: SplitRatioSelectorProps) {
  function splitTypeLabel(splitType: string) {
    switch (splitType) {
      case 'equal':
        return 'igualmente'
      case 'proportional':
        return 'proporcionalmente'
      case 'transfer':
        return 'transferência'
    }
  }

  return (
    <FormField
      name="splitType"
      control={form.control}
      render={({ field }) => (
        <FormItem>
          <Select onValueChange={field.onChange} defaultValue={field.value}>
            <FormControl>
              <CustomSelectTrigger>
                <Badge>{splitTypeLabel(field.value)}</Badge>
              </CustomSelectTrigger>
            </FormControl>
            <SelectContent>
              <SelectItem value="equal">igualmente</SelectItem>
              <SelectItem value="proportional">proporcionalmente</SelectItem>
              <SelectItem value="transfer">transferência</SelectItem>
            </SelectContent>
          </Select>
        </FormItem>
      )}
    />
  )
}

export { SplitRatioSelector }

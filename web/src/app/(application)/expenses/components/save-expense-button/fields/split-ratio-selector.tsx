'use client'

import { useTranslations } from 'next-intl'
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
  const t = useTranslations()

  function splitTypeLabel(splitType: string) {
    switch (splitType) {
      case 'equal':
        return t('expenses.splitType.equalLabel')
      case 'proportional':
        return t('expenses.splitType.proportionalLabel')
      case 'transfer':
        return t('expenses.splitType.transferLabel')
      default:
        return ''
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
              <SelectItem value="equal">{t('expenses.splitType.equalLabel')}</SelectItem>
              <SelectItem value="proportional">{t('expenses.splitType.proportionalLabel')}</SelectItem>
              <SelectItem value="transfer">{t('expenses.splitType.transferLabel')}</SelectItem>
            </SelectContent>
          </Select>
        </FormItem>
      )}
    />
  )
}

export { SplitRatioSelector }

'use client'

import { useTranslations } from 'next-intl'
import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { expenseFormSchema } from '@/schemas'

type DescriptionFieldProps = {
  form: UseFormReturn<z.infer<typeof expenseFormSchema>>
}

function DescriptionField({ form }: DescriptionFieldProps) {
  const t = useTranslations()

  return (
    <FormField
      name="description"
      control={form.control}
      render={({ field }) => (
        <FormItem>
          <FormControl>
            <Input
              autoFocus
              placeholder={t('expenses.descriptionPlaceholder')}
              {...field}
              variant="outline"
              autoCapitalize="words"
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  )
}

export { DescriptionField }

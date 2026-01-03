import React from 'react'
import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { FormField, FormItem } from '@/components/ui/form'
import { Spinning } from '@/components/ui/spinning'
import { Group } from '@/domain/group'
import { User } from '@/domain/user'
import { useCategories } from '@/hooks/use-categories'
import { expenseFormSchema } from '@/schemas'

import { CategorySelectorMenu } from '../../category_selector'

type CategorySelectorProps = {
  form: UseFormReturn<z.infer<typeof expenseFormSchema>>
  user?: User
  group?: Group
  className?: string
  iconSize?: number
  isPredicting?: boolean
  isPredicted?: boolean
  onCategorySelect?: (value: number) => void
}

function CategorySelector({
  form,
  className,
  iconSize,
  isPredicting,
  isPredicted,
  onCategorySelect,
}: CategorySelectorProps) {
  const [open, setOpen] = React.useState(false)
  const { categories, getCategory, isLoading } = useCategories()

  function handleClick(value: number) {
    form.setValue('categoryId', value)
    setOpen(false)
    onCategorySelect?.(value)
  }

  if (isLoading) {
    return (
      <div className="mr-3 flex items-center justify-center">
        <Spinning />
      </div>
    )
  }

  return (
    <FormField
      name="categoryId"
      control={form.control}
      render={({ field }) => (
        <FormItem className={className}>
          <CategorySelectorMenu
            open={open}
            setOpen={setOpen}
            handleClick={handleClick}
            getCategory={getCategory}
            selectedICategoryID={field.value}
            iconSize={iconSize}
            isPredicting={isPredicting}
            isPredicted={isPredicted}
            categories={categories}
          />
        </FormItem>
      )}
    />
  )
}

export { CategorySelector }

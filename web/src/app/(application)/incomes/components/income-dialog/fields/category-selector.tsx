import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { FormControl, FormField, FormItem } from '@/components/ui/form'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Group } from '@/domain/group'
import { User } from '@/domain/user'
import { incomeFormSchema } from '@/schemas'

type CategorySelectorProps = {
  form: UseFormReturn<z.infer<typeof incomeFormSchema>>
  user?: User
  group?: Group
}

function CategorySelector({ form }: CategorySelectorProps) {
  return (
    <FormField
      name="type"
      control={form.control}
      render={({ field }) => (
        <FormItem>
          <Select onValueChange={field.onChange} defaultValue={field.value}>
            <FormControl>
              <SelectTrigger className="w-[240px]">
                <SelectValue placeholder="Tipo de receita" />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="salary">Salário</SelectItem>
                <SelectItem value="benefit">Benefício</SelectItem>
                <SelectItem value="vacation">Férias</SelectItem>
                <SelectItem value="thirteenth_salary">Décimo terceiro</SelectItem>
                <SelectItem value="other">Outro</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormItem>
      )}
    />
  )
}

export { CategorySelector }

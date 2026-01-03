import { UseFormReturn } from 'react-hook-form'
import { z } from 'zod'

import { FormControl, FormField, FormItem } from '@/components/ui/form'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Group } from '@/domain/group'
import { incomeFormSchema } from '@/schemas'

type CategorySelectorProps = {
  form: UseFormReturn<z.infer<typeof incomeFormSchema>>
  group?: Group
}

function UserSelector({ form, group }: CategorySelectorProps) {
  return (
    <FormField
      name="userId"
      control={form.control}
      render={({ field }) => (
        <FormItem>
          <Select onValueChange={(value) => field.onChange(parseInt(value))} defaultValue={String(field.value)}>
            <FormControl>
              <SelectTrigger className="w-[240px]">
                <SelectValue placeholder="Tipo de receita" />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              <SelectGroup>
                {group?.members.map((member) => (
                  <SelectItem key={member.id} value={String(member.id)}>
                    {member.name}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormItem>
      )}
    />
  )
}

export { UserSelector }

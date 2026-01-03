import { format } from 'date-fns'
import { CalendarIcon } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { TableCell, TableRow } from '@/components/ui/table'
import { useCategories } from '@/hooks/use-categories'
import { useGroup } from '@/hooks/use-group'
import { moneyMask } from '@/lib/utils'
import { ExpenseDraftRow } from '@/types/expense-import'

import { CategorySelectorMenu } from './category_selector'

interface ExpenseDraftRowProps {
  draft: ExpenseDraftRow
  onUpdate: (draftId: string, updater: Partial<ExpenseDraftRow>) => void
  isPredicting: boolean
  isPredicted: boolean
}

export function ExpenseDraftRowComponent({ draft, onUpdate, isPredicting, isPredicted }: ExpenseDraftRowProps) {
  const { me, partner } = useGroup()
  const { categories, getCategory } = useCategories()

  return (
    <TableRow
      className={`align-top ${
        draft.status === 'success'
          ? 'bg-green-50 dark:bg-green-950/20'
          : draft.status === 'error'
          ? 'bg-red-50 dark:bg-red-950/20'
          : draft.status === 'saving'
          ? 'bg-blue-50 opacity-60 dark:bg-blue-950/20'
          : ''
      }`}
    >
      <TableCell className="text-center">
        <Checkbox
          checked={draft.include}
          onCheckedChange={(checked) => onUpdate(draft.id, { include: checked === true })}
        />
      </TableCell>
      <TableCell>
        <Input
          value={draft.description}
          onChange={(event) => onUpdate(draft.id, { description: event.currentTarget.value })}
          placeholder="Descrição"
          className="min-w-[200px]"
        />
      </TableCell>
      <TableCell>
        <Popover>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              className={`w-full justify-start text-left font-normal ${!draft.date && 'text-muted-foreground'}`}
            >
              <CalendarIcon className="mr-2 h-4 w-4" />
              {draft.date ? format(new Date(draft.date), 'dd/MM/yyyy') : 'Selecione'}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto p-0" align="start">
            <Calendar
              mode="single"
              selected={draft.date ? new Date(draft.date) : undefined}
              onSelect={(date) =>
                onUpdate(draft.id, {
                  date: date?.toISOString(),
                })
              }
              disabled={(date) => date < new Date('1900-01-01') || date > new Date()}
            />
          </PopoverContent>
        </Popover>
      </TableCell>
      <TableCell>
        <Input
          placeholder="R$ 0,00"
          inputMode="decimal"
          prefix="R$"
          onChange={(event) => onUpdate(draft.id, { amountInCents: event.currentTarget.value.replace(/\D/g, '') })}
          value={moneyMask(draft.amountInCents)}
          className="min-w-[120px]"
        />
      </TableCell>
      <TableCell>
        <CategorySelectorMenu
          handleClick={(categoryId) => onUpdate(draft.id, { categoryId })}
          getCategory={getCategory}
          selectedICategoryID={draft.categoryId}
          iconSize={20}
          isPredicting={isPredicting}
          isPredicted={isPredicted}
          categories={categories}
        />
      </TableCell>
      <TableCell>
        <Select
          value={draft.payerId ? String(draft.payerId) : undefined}
          onValueChange={(value) => onUpdate(draft.id, { payerId: Number(value) })}
        >
          <SelectTrigger>
            <SelectValue placeholder="Selecione" />
          </SelectTrigger>
          <SelectContent>
            {me && <SelectItem value={String(me.id)}>Você</SelectItem>}
            {partner && <SelectItem value={String(partner.id)}>{partner.name}</SelectItem>}
          </SelectContent>
        </Select>
      </TableCell>
      <TableCell>
        <Select
          value={draft.splitType}
          onValueChange={(value: 'equal' | 'proportional' | 'transfer') => onUpdate(draft.id, { splitType: value })}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="equal">Igualmente</SelectItem>
            <SelectItem value="proportional">Proporcionalmente</SelectItem>
            <SelectItem value="transfer">Transferência</SelectItem>
          </SelectContent>
        </Select>
      </TableCell>
    </TableRow>
  )
}

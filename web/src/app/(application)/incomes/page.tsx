'use client'

import { PencilIcon, Trash2Icon } from 'lucide-react'
import { useState } from 'react'

import { IncomeDialog } from '@/app/(application)/incomes/components/income-dialog'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableFooter, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { YearMonthPicker } from '@/components/year-month-picker'
import { incomeLabel } from '@/domain/income'
import { useGroup } from '@/hooks/use-group'
import { useIncome } from '@/hooks/use-income'
import { formatCurrency } from '@/lib/utils'

import { DeleteIncomeButton } from './components/delete-income-button'

export default function Incomes() {
  const [date, setDate] = useState(new Date())

  const { getMember, me } = useGroup()
  const { incomes, totalIncome, isLoading } = useIncome(date)

  return (
    <div className="mx-auto flex flex-col gap-6 p-2 md:min-w-[56rem]">
      <div className="flex justify-between">
        <YearMonthPicker selectedDate={date} onSelectDate={setDate} />
        <IncomeDialog type="create" date={date}>
          <Button>Nova receita</Button>
        </IncomeDialog>
      </div>
      <Table className="w-full">
        <TableHeader>
          <TableRow>
            <TableHead className="w-[200px]">Nome</TableHead>
            <TableHead className="w-[200px]">Tipo</TableHead>
            <TableHead className="w-[200px]">Valor</TableHead>
            <TableHead className="w-[100px]">Ação</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading && <LoadingTable />}
          {incomes.map((income) => (
            <TableRow key={income.id}>
              <TableCell>{getMember(income.userId)?.name}</TableCell>
              <TableCell>{incomeLabel(income.type)}</TableCell>
              <TableCell>{formatCurrency(income.amount)}</TableCell>
              <TableCell>
                {(me?.flags?.find((f) => f === 'edit_partner_income') || income.userId === me?.id) && (
                  <>
                    <IncomeDialog income={income} type="update" date={date}>
                      <Button size="icon" variant="ghost">
                        <PencilIcon className="h-4 w-4" />
                      </Button>
                    </IncomeDialog>
                    <DeleteIncomeButton income={income}>
                      <Button size="icon" variant="ghost">
                        <Trash2Icon className="h-4 w-4" />
                      </Button>
                    </DeleteIncomeButton>
                  </>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
        {incomes.length > 0 && (
          <TableFooter>
            <TableRow>
              <TableCell className="font-bold" colSpan={2}>
                Total
              </TableCell>
              <TableCell className="font-bold">{formatCurrency(totalIncome)}</TableCell>
              <TableCell />
            </TableRow>
          </TableFooter>
        )}
      </Table>
    </div>
  )
}

function LoadingTable() {
  return Array.from({ length: 3 }).map((_, i) => (
    <TableRow key={i}>
      {Array.from({ length: 4 }).map((_, j) => (
        <TableCell key={j}>
          <Skeleton className="h-4 w-12" />
        </TableCell>
      ))}
    </TableRow>
  ))
}

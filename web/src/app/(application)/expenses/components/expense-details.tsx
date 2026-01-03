/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import React from 'react'

import { Spinning } from '@/components/ui/spinning'
import { Category } from '@/domain/category'
import { Expense } from '@/domain/expense'
import { User } from '@/domain/user'
import { useCategories } from '@/hooks/use-categories'
import { useExpenseDetails } from '@/hooks/use-expense-details'
import { useGroup } from '@/hooks/use-group'
import { cn, findDifferences, formatCurrency } from '@/lib/utils'

type ExpenseDetailsProps = React.HTMLAttributes<HTMLDivElement> & {
  expense: Expense
}

function ExpenseDetails({ expense, className, ...props }: ExpenseDetailsProps) {
  const { data, isLoading } = useExpenseDetails(expense.id)
  const { getMember } = useGroup()
  const { getCategory } = useCategories()
  return (
    <div className={cn('flex flex-col', className)} {...props}>
      <span className="text-sm font-semibold">Histórico:</span>
      {isLoading ? (
        <Spinning />
      ) : (
        <ul className="list-disc pl-4">
          {data?.map((expense, i) => {
            if (i === 0) {
              return (
                <li key={i}>
                  <div className="flex justify-between">
                    <span>{`Despesa '${expense.name}' criada no valor de '${formatCurrency(expense.amount)}'`}</span>
                    <span className="ml-1">{format(expense.updatedAt!, 'dd/MM/yyyy')}</span>
                  </div>
                </li>
              )
            }

            const descriptions = translateDiff(findDifferences(data[i - 1], expense), getMember, getCategory)
            return (
              <li key={i}>
                <div className="flex justify-between">
                  <div className="flex flex-col">
                    {descriptions.map((description, i) => (
                      <span key={i}>{description}</span>
                    ))}
                  </div>
                  <span className="ml-1">{format(expense.updatedAt!, 'dd/MM/yyyy', { locale: ptBR })}</span>
                </div>
              </li>
            )
          })}
        </ul>
      )}
    </div>
  )
}

function translateDiff(
  record: Record<string, { old: string; actual: string }>,
  getMember: (userId: number) => User | undefined,
  getCategory: (categoryId: number) => Category | undefined,
) {
  function translateBuilder(phrase: string, old?: string, actual?: string) {
    return phrase.replace('old', old ?? '').replace('actual', actual ?? '')
  }

  const phraseByKey = {
    name: "Nome alterado de 'old' para 'actual'",
    amount: "Valor alterado de 'old' para 'actual'",
    refundAmount: "Reembolso alterado de 'old' para 'actual'",
    description: "Descrição alterado de 'old' para 'actual'",
    categoryId: "Categoria alterada de 'old' para 'actual'",
    payerId: "Pagador alterado de 'old' para 'actual'",
    receiverId: "Recebedor alterado de 'old' para 'actual'",
    'splitRatio.payer': "Proporção do pagador alterada de 'old' para 'actual'",
    'splitRatio.receiver': "Proporção do recebedor alterada de 'old' para 'actual'",
    createdAt: "Data de criação alterada de 'old' para 'actual'",
  } as Record<string, string>

  const recordKeys = Object.keys(record)
  const allowedKeys = Object.keys(phraseByKey)
  const filteredKeys = allowedKeys.filter((key) => recordKeys.includes(key))

  return filteredKeys.map((key) => {
    if (key === 'refundAmount' && record[key].old === null) {
      return translateBuilder(
        "Adição de reembolso no valor de 'actual'",
        record[key].old,
        formatCurrency(parseInt(record[key].actual)),
      )
    }

    if (key === 'amount' || key === 'refundAmount') {
      return translateBuilder(
        phraseByKey[key],
        formatCurrency(parseInt(record[key].old)),
        formatCurrency(parseInt(record[key].actual)),
      )
    }

    if (key === 'categoryId') {
      return translateBuilder(
        phraseByKey[key],
        getCategory(parseInt(record[key].old))?.name,
        getCategory(parseInt(record[key].actual))?.name,
      )
    }

    if (key === 'payerId' || key === 'receiverId') {
      return translateBuilder(
        phraseByKey[key],
        getMember(parseInt(record[key].old))?.name,
        getMember(parseInt(record[key].actual))?.name,
      )
    }

    if (key === 'createdAt') {
      return translateBuilder(
        phraseByKey[key],
        format(new Date(record[key].old), 'dd/MM/yyyy'),
        format(new Date(record[key].actual), 'dd/MM/yyyy'),
      )
    }

    return translateBuilder(phraseByKey[key], record[key].old, record[key].actual)
  })
}

export { ExpenseDetails }

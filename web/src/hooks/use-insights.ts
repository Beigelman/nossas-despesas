import { useQuery } from '@tanstack/react-query'
import { DateRange } from 'react-day-picker'

import { privateHttpClient } from '@/http/private-client'
import { ApiResponse } from '@/lib/api'
import { newUTCDate } from '@/lib/date'

type GetExpensePerCategoryApiResponse = ApiResponse<
  {
    name: string
    amount: number
    categories: {
      name: string
      amount: number
    }[]
  }[]
>

type GetExpensePerPeriodApiResponse = ApiResponse<
  {
    date: string
    amount: number
    quantity: number
  }[]
>

type UseInsightsProps = {
  range: DateRange
  aggregate: 'month' | 'day'
}

function useInsights({ range, aggregate }: UseInsightsProps) {
  const from = newUTCDate(range?.from)
  const to = newUTCDate(range?.to)

  const { data: incomesPerPeriod, isLoading: incomesPerPeriodIsLoading } = useQuery({
    queryKey: ['insights-incomes-period', from, to, aggregate],
    queryFn: async () => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetExpensePerPeriodApiResponse>('/incomes/insights', {
        params: { start_date: from, end_date: to, aggregate },
      })

      const { total, count } = data.reduce(
        (acc, period) => ({ total: acc.total + period.amount, count: acc.count + period.quantity }),
        { total: 0, count: 0 },
      )

      return {
        data,
        count,
        total,
      }
    },
  })

  const { data: expensesPerPeriod, isLoading: expensesPerPeriodIsLoading } = useQuery({
    queryKey: ['insights-expenses-period', from, to, aggregate],
    queryFn: async () => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetExpensePerPeriodApiResponse>('/expenses/insights', {
        params: { start_date: from, end_date: to, aggregate },
      })

      const { total, count } = data.reduce(
        (acc, period) => ({ total: acc.total + period.amount, count: acc.count + period.quantity }),
        { total: 0, count: 0 },
      )

      return {
        data,
        count,
        total,
      }
    },
  })

  const { data: expensesPerCategory, isLoading: expensesPerCategoryIsLoading } = useQuery({
    queryKey: ['insights-expenses-categories', from, to],
    queryFn: async () => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetExpensePerCategoryApiResponse>('/expenses/insights/category', {
        params: { start_date: from, end_date: to },
      })

      return data
        .sort((a, b) => b.amount - a.amount)
        .map((category) => ({
          ...category,
          categories: category.categories.sort((a, b) => b.amount - a.amount),
        }))
    },
  })

  return {
    expensesPerCategory,
    expensesPerCategoryIsLoading,
    expensesPerPeriod,
    expensesPerPeriodIsLoading,
    incomesPerPeriod,
    incomesPerPeriodIsLoading,
  }
}

export { useInsights }

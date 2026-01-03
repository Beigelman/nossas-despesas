import { ArrowBigDownDash, ArrowBigUpDash, ShoppingBag } from 'lucide-react'
import { DateRange } from 'react-day-picker'
import { Tooltip } from 'recharts'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useInsights } from '@/hooks/use-insights'
import { formatCurrency } from '@/lib/utils'

import { BarChart } from './bar-chart'

type ExpensesPeriodInsightProps = {
  dateRange: DateRange
  aggregate: 'month' | 'day'
}

function ExpensesPeriodInsight({ dateRange, aggregate }: ExpensesPeriodInsightProps) {
  const {
    expensesPerCategory,
    expensesPerCategoryIsLoading,
    expensesPerPeriod,
    expensesPerPeriodIsLoading,
    incomesPerPeriod,
    incomesPerPeriodIsLoading,
  } = useInsights({
    range: dateRange,
    aggregate,
  })

  return (
    <>
      <div className="grid gap-4 md:grid-cols-1 lg:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Gastos</CardTitle>
            <ArrowBigUpDash color="#888888" size={18} />
          </CardHeader>
          <CardContent>
            {expensesPerPeriodIsLoading ? (
              <Skeleton className="h-7 w-12" />
            ) : (
              <div className="text-2xl font-bold">{formatCurrency(expensesPerPeriod?.total)}</div>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Receita</CardTitle>
            <ArrowBigDownDash color="#888888" size={18} />
          </CardHeader>
          <CardContent>
            {incomesPerPeriodIsLoading ? (
              <Skeleton className="h-7 w-12" />
            ) : (
              <div className="text-2xl font-bold">{formatCurrency(incomesPerPeriod?.total)}</div>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Número de compras</CardTitle>
            <ShoppingBag color="#888888" size={18} />
          </CardHeader>
          <CardContent>
            {expensesPerPeriodIsLoading ? (
              <Skeleton className="h-5 w-10" />
            ) : (
              <div className="text-2xl font-bold">{expensesPerPeriod?.count}</div>
            )}
          </CardContent>
        </Card>
      </div>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-7">
        <Card className="col-span-1 md:col-span-4">
          <CardHeader>
            <CardTitle>Overview</CardTitle>
          </CardHeader>
          <CardContent className="pl-2">
            <BarChart
              data={expensesPerPeriod?.data ?? []}
              XKey="date"
              YKey="amount"
              direction="horizontal"
              isLoading={expensesPerPeriodIsLoading}
              customTooltip={
                <Tooltip
                  content={({ payload }) =>
                    payload?.map((entry, i) => (
                      <div key={i} className="flex flex-col rounded-md border bg-card p-3 shadow-sm dark:bg-zinc-900">
                        <span>{`Data: ${entry.payload?.date}`}</span>
                        <span>{`Total: ${formatCurrency(entry.payload?.amount)}`}</span>
                      </div>
                    ))
                  }
                />
              }
            />
          </CardContent>
        </Card>
        <Card className="col-span-1 md:col-span-3">
          <CardHeader>
            <CardTitle>Categorias</CardTitle>
            <CardDescription>Visão por grupo de categorias</CardDescription>
          </CardHeader>
          <CardContent>
            <BarChart
              data={expensesPerCategory}
              XKey="name"
              YKey="amount"
              direction="vertical"
              isLoading={expensesPerCategoryIsLoading}
              customTooltip={
                <Tooltip
                  content={({ payload }) =>
                    payload?.map((entry, i) => (
                      <div key={i} className="flex flex-col rounded-md border bg-card p-3 shadow-sm dark:bg-zinc-900">
                        <div className="mb-2 flex justify-between gap-3">
                          <span>Total</span>
                          <span>{formatCurrency(entry.payload?.amount)}</span>
                        </div>
                        {entry.payload?.categories?.map((c: { name: string; amount: number }) => (
                          <div key={c.name} className="flex justify-between gap-3">
                            <span>{c.name}:</span>
                            <span>{formatCurrency(c.amount)}</span>
                          </div>
                        ))}
                      </div>
                    ))
                  }
                />
              }
            />
          </CardContent>
        </Card>
      </div>
    </>
  )
}

export { ExpensesPeriodInsight }

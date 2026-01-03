import { Fragment, useEffect, useState } from 'react'
import { useInView } from 'react-intersection-observer'

import { ExpenseLine } from '@/app/(application)/expenses/components/expense-line'
import { MonthSeparator } from '@/components/month-separator'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinning } from '@/components/ui/spinning'
import { useExpenses } from '@/hooks/use-expenses'
import { useSearch } from '@/hooks/use-search'
import { cn } from '@/lib/utils'

function ExpensesList() {
  const { ref, inView } = useInView()
  const search = useSearch((state) => state.search)
  const { expensesData } = useExpenses(search)
  const { data, isLoading, fetchNextPage, hasNextPage, isFetchingNextPage, isFetching } = expensesData
  const [selectedExpense, setSelectedExpense] = useState(0)

  useEffect(() => {
    if (inView) {
      fetchNextPage()
    }
  }, [fetchNextPage, inView])

  if (isLoading) {
    return (
      <div className="flex h-full flex-col gap-2 p-3">
        {Array.from({ length: 10 }).map((_, i) => (
          <div key={i} className="flex items-center justify-between space-x-4">
            <Skeleton className="h-12 w-12 rounded-md" />
            <Skeleton className="h-4 w-[200px]" />
            <div className="flex gap-2">
              {Array.from({ length: 2 }).map((_, j) => (
                <div key={j} className="flex flex-col items-end gap-1">
                  <Skeleton className="h-4 w-[100px]" />
                  <Skeleton className="h-4 w-[70px]" />
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (!isLoading && (data?.pages.length === 0 || data?.pages[0].expenses.length === 0)) {
    return (
      <div className="flex h-full flex-col items-center justify-center">
        <p className="text-center">Nenhum registro encontrado</p>
        <p className="text-center text-sm">Comece adicionando um novo gasto</p>
      </div>
    )
  }

  return (
    <Card>
      <CardContent className="overflow-hidden p-0">
        {data?.pages.map((page, i) =>
          page.expenses.map((expense, j) => {
            const currentDate = expense.createdAt
            const nextDate = page.expenses[j + 1]
              ? page.expenses[j + 1]?.createdAt
              : data.pages[i + 1]?.expenses[0]?.createdAt
            return (
              <Fragment key={expense.id}>
                <ExpenseLine
                  expense={expense}
                  selectedExpense={selectedExpense}
                  setSelectedExpense={setSelectedExpense}
                />
                <MonthSeparator nextDate={nextDate} currentDate={currentDate} />
              </Fragment>
            )
          }),
        )}
        <Button
          ref={ref}
          onClick={() => fetchNextPage()}
          disabled={!hasNextPage || isFetchingNextPage}
          variant={'ghost'}
          className={cn('w-full', hasNextPage ? 'block' : 'hidden')}
        >
          {isFetchingNextPage ? <Spinning className="h-5 w-5" /> : hasNextPage ? 'Carregar mais items' : ''}
        </Button>
        <div>{isFetching && !isFetchingNextPage ? 'Background Updating...' : null}</div>
      </CardContent>
    </Card>
  )
}

export { ExpensesList }

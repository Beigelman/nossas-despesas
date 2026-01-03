'use client'

import { PlusIcon } from 'lucide-react'

import { ImportExpensesButton } from '@/app/(application)/expenses/components/import-exepenses-button'
import { Overview } from '@/app/(application)/expenses/components/overview'
import { SaveExpenseButton } from '@/app/(application)/expenses/components/save-expense-button'
import { Button } from '@/components/ui/button'

import { ExpensesList } from './components/expenses-list'

export default function Expenses() {
  return (
    <div className="mx-auto flex w-screen flex-col-reverse md:w-max md:flex-row">
      <div className="flex h-full w-full flex-col gap-2 py-3 md:w-[500px] lg:w-[750px]">
        <div className="flex items-center justify-between gap-3">
          <ImportExpensesButton />
          <SaveExpenseButton type="create">
            <Button
              variant="default"
              className="group fixed bottom-6 right-6 z-50 h-12 w-min gap-1 overflow-hidden rounded-full md:static md:z-0 md:mb-3 md:ml-auto md:mr-3 md:h-10"
            >
              <PlusIcon size={18} />
              <span
                className={`hidden w-0 translate-x-32 whitespace-nowrap transition-all duration-500 group-hover:w-32 group-hover:translate-x-0 md:block`}
              >
                Adicionar despesa
              </span>
            </Button>
          </SaveExpenseButton>
        </div>
        <Overview />
        <ExpensesList />
      </div>
    </div>
  )
}

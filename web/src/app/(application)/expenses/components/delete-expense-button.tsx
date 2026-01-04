'use client'

import { useTranslations } from 'next-intl'
import { useState } from 'react'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer'
import { Expense } from '@/domain/expense'
import { useExpenses } from '@/hooks/use-expenses'
import { useMediaQuery } from '@/hooks/use-media-query'
type DeleteExpenseButtonProps = {
  children: React.ReactNode
  expense: Expense
}

function DeleteExpenseButton({ expense, children }: DeleteExpenseButtonProps) {
  const [open, setOpen] = useState(false)
  const isDesktop = useMediaQuery('(min-width: 768px)')
  const { deleteExpense } = useExpenses('')
  const t = useTranslations()

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>{children}</DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t('expenses.confirmDelete')}</DialogTitle>
            <DialogDescription>{`${t('expenses.deleteExpense')}: ${expense.name}`}</DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2">
            <DialogClose>
              <Button variant="secondary" className="w-full md:w-fit">
                {t('common.cancel')}
              </Button>
            </DialogClose>
            <Button variant="destructive" onClick={() => deleteExpense(expense.id)}>
              {t('common.delete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    )
  }

  return (
    <Drawer open={open} onOpenChange={setOpen} shouldScaleBackground={false}>
      <DrawerTrigger asChild>{children}</DrawerTrigger>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>{t('expenses.confirmDelete')}</DrawerTitle>
          <DrawerDescription>{`${t('expenses.deleteExpense')}: ${expense.name}`}</DrawerDescription>
        </DrawerHeader>
        <DrawerFooter>
          <DrawerClose>
            <Button variant="secondary" className="w-full md:w-fit">
              {t('common.cancel')}
            </Button>
          </DrawerClose>
          <Button variant="destructive" onClick={() => deleteExpense(expense.id)}>
            {t('common.delete')}
          </Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}

export { DeleteExpenseButton }

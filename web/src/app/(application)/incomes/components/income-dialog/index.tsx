'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo, useState } from 'react'
import { useForm } from 'react-hook-form'
import * as z from 'zod'

import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer'
import { Income } from '@/domain/income'
import { useGroup } from '@/hooks/use-group'
import { useIncome } from '@/hooks/use-income'
import { useMediaQuery } from '@/hooks/use-media-query'
import { incomeFormSchema } from '@/schemas'

import { Button } from '../../../../../components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../../../../../components/ui/dialog'
import { Form } from '../../../../../components/ui/form'
import { AmountField } from './fields/amout-field'
import { CategorySelector } from './fields/category-selector'
import { DateSelector } from './fields/date-selector'
import { UserSelector } from './fields/user-selector'

type IncomeDialogProps = {
  type: 'create' | 'update'
  children?: React.ReactNode
  income?: Income
  date?: Date
}

function IncomeDialog({ children, income, type, date }: IncomeDialogProps) {
  const [open, setOpen] = useState(false)

  const isDesktop = useMediaQuery('(min-width: 768px)')
  const { createIncome, updateIncome } = useIncome()

  const { group, me } = useGroup()

  const form = useForm<z.infer<typeof incomeFormSchema>>({
    resolver: zodResolver(incomeFormSchema),
    defaultValues: {
      amount: '',
      date: new Date(),
      type: 'salary',
      userId: me?.id,
    },
  })

  const title = useMemo(() => (type === 'create' ? 'Adicionar nova receita' : 'Atualizar receita'), [type])

  useEffect(() => {
    form.setValue('amount', income?.amount ? `${income?.amount}` : '')
    form.setValue('date', income?.createdAt ?? date ?? new Date())
    form.setValue('type', income?.type ?? 'salary')
    form.setValue('userId', income?.userId ?? me?.id ?? 0)
  }, [form, income, date, me])

  function onSubmit({ amount, date, type: incomeType, userId }: z.infer<typeof incomeFormSchema>) {
    const payload = {
      amount: parseInt(amount),
      type: incomeType,
      date,
      user_id: userId,
    }

    try {
      if (type === 'update') {
        updateIncome({ ...payload, id: income?.id ?? 0 })
      } else {
        createIncome(payload)
      }
    } catch (e) {
      console.error(e)
    } finally {
      form.reset()
      setOpen(false)
    }
  }

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>{children}</DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
          </DialogHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <div className="flex flex-col items-center gap-2">
                <AmountField form={form} />
                <UserSelector form={form} group={group} />
                <CategorySelector form={form} />
                <DateSelector form={form} date={date} />
              </div>
              <DialogFooter>
                <div className="mt-4 flex flex-row-reverse gap-2">
                  <Button type="submit">Salvar</Button>
                  <DialogClose>
                    <Button type="button" variant="secondary" className="w-full md:w-fit">
                      Cancelar
                    </Button>
                  </DialogClose>
                </div>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    )
  }

  return (
    <Drawer open={open} onOpenChange={setOpen} shouldScaleBackground={false}>
      <DrawerTrigger asChild>{children}</DrawerTrigger>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>{title}</DrawerTitle>
        </DrawerHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="flex flex-col items-center gap-2">
              <AmountField form={form} />
              <UserSelector form={form} group={group} />
              <CategorySelector form={form} />
              <DateSelector form={form} date={date} />
            </div>
            <DrawerFooter className="mt-4 gap-2">
              <DrawerClose>
                <Button type="button" variant="secondary" className="w-full md:w-fit">
                  Cancelar
                </Button>
              </DrawerClose>
              <Button type="submit">Salvar</Button>
            </DrawerFooter>
          </form>
        </Form>
      </DrawerContent>
    </Drawer>
  )
}

export { IncomeDialog }

'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo, useRef, useState } from 'react'
import { useForm } from 'react-hook-form'
import * as z from 'zod'

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
import { useGroup } from '@/hooks/use-group'
import { useMediaQuery } from '@/hooks/use-media-query'
import { usePredict } from '@/hooks/use-predict'
import { isSameDate } from '@/lib/date'
import { expenseFormSchema } from '@/schemas'

import { Button } from '../../../../../components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../../../../../components/ui/dialog'
import { Form } from '../../../../../components/ui/form'
import { AmountField } from './fields/amout-field'
import { CategorySelector } from './fields/category-selector'
import { DateSelector } from './fields/date-selector'
import { DescriptionField } from './fields/description-field'
import { PayerSelector } from './fields/payer-selector'
import { RefundField } from './fields/refund-field'
import { SplitRatioSelector } from './fields/split-ratio-selector'

type SaveExpenseButtonProps = {
  children?: React.ReactNode
  expense?: Expense
  type: 'create' | 'update'
}

function SaveExpenseButton({ children, expense, type }: SaveExpenseButtonProps) {
  const [open, setOpen] = useState(false)

  const isDesktop = useMediaQuery('(min-width: 768px)')
  const { group, me, partner } = useGroup()
  const { createExpense, updateExpense } = useExpenses('')
  const { predictCategoryID, isPredicting } = usePredict()

  const form = useForm<z.infer<typeof expenseFormSchema>>({
    resolver: zodResolver(expenseFormSchema),
    defaultValues: {
      description: '',
      date: new Date(),
      amount: '',
      payerId: me?.id,
      splitType: 'proportional',
      categoryId: 16, // Supermarket category
      refundAmount: '',
    },
  })

  const [wasPredicted, setWasPredicted] = useState(false)
  const lastPredictionInputRef = useRef<{ description: string; amount: string } | null>(null)
  const userSelectedCategoryRef = useRef(false)

  const watchedDescription = form.watch('description')
  const watchedAmount = form.watch('amount')

  const title = useMemo(() => (type === 'create' ? 'Adicionar uma nova despesa' : 'Atualizar despesa'), [type])

  const description = useMemo(
    () => (type === 'create' ? `No grupo ${group?.name}` : `${expense?.name}`),
    [expense, group, type],
  )

  useEffect(() => {
    form.setValue('description', expense?.name ?? '')
    form.setValue('amount', expense?.amount ? `${expense?.amount}` : '')
    form.setValue('payerId', expense?.payerId ?? me?.id ?? 0)
    form.setValue('date', expense?.createdAt ?? new Date())
    form.setValue('splitType', expense?.splitType ?? 'proportional')
    form.setValue('categoryId', expense?.categoryId ?? 16) // Supermarket category
    form.setValue('refundAmount', expense?.refundAmount ? `${expense?.refundAmount}` : '')
  }, [form, me?.id, expense])

  useEffect(() => {
    const trimmedDescription = watchedDescription ?? ''
    const normalizedAmount = watchedAmount ?? ''

    const hasEnoughDescription = trimmedDescription.length >= 3
    const hasAmount = normalizedAmount.length > 0

    if (!hasEnoughDescription && !hasAmount) {
      setWasPredicted(false)
      lastPredictionInputRef.current = null
      return
    }

    if (
      lastPredictionInputRef.current &&
      lastPredictionInputRef?.current?.description === trimmedDescription &&
      lastPredictionInputRef?.current?.amount === normalizedAmount
    ) {
      return
    }

    setWasPredicted(false)
    userSelectedCategoryRef.current = false

    const handler = setTimeout(() => {
      predictCategoryID({
        name: trimmedDescription,
        amount_cents: Number(normalizedAmount) || 0,
      })
        .then((response) => {
          const predictedCategoryId = response.data.category_id

          if (typeof predictedCategoryId !== 'number') {
            return
          }

          if (userSelectedCategoryRef.current) {
            return
          }

          lastPredictionInputRef.current = {
            description: trimmedDescription,
            amount: normalizedAmount,
          }

          form.setValue('categoryId', predictedCategoryId)
          setWasPredicted(true)
        })
        .catch((error) => {
          console.error(error)
        })
    }, 600)

    return () => clearTimeout(handler)
  }, [watchedAmount, watchedDescription, form, predictCategoryID])

  function onSubmit({
    amount,
    categoryId,
    description,
    payerId,
    splitType,
    date,
    refundAmount,
  }: z.infer<typeof expenseFormSchema>) {
    const receiverId = payerId === me?.id ? partner?.id : me?.id
    const payload = {
      amount: parseInt(amount),
      name: description,
      receiverId: receiverId ?? 0,
      createdAt: date,
      payerId,
      splitType,
      categoryId,
      refundAmount: refundAmount !== '' ? parseInt(refundAmount) : 0,
    }

    try {
      if (type === 'update' && expense) {
        updateExpense({
          ...payload,
          createdAt: isSameDate(expense.createdAt, date) ? undefined : date,
          id: expense.id,
        })
      } else {
        createExpense(payload)
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
            <DialogDescription>{description}</DialogDescription>
          </DialogHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <div className="flex flex-col gap-4 px-10">
                <div className="flex w-full flex-row-reverse items-center justify-center gap-3">
                  <div className="flex w-full flex-col">
                    <DescriptionField form={form} />
                    <AmountField form={form} />
                  </div>
                  <CategorySelector
                    form={form}
                    className="mt-5 flex"
                    isPredicting={isPredicting}
                    isPredicted={wasPredicted}
                    onCategorySelect={() => {
                      userSelectedCategoryRef.current = true
                      setWasPredicted(false)
                    }}
                  />
                </div>
                <div className="mx-auto flex items-center space-x-1">
                  <span className="mt-1 text-sm">Pago por</span>
                  <PayerSelector form={form} user={me} partner={partner} />
                  <span className="mt-1 text-sm">dividido</span>
                  <SplitRatioSelector form={form} />
                </div>
                <DateSelector form={form} />
                {type === 'update' && <RefundField form={form} />}
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
          <DrawerDescription>{description}</DrawerDescription>
        </DrawerHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="flex flex-col gap-4 px-10">
              <div className="flex w-full flex-row-reverse items-center justify-center gap-3">
                <div className="flex w-full flex-col">
                  <DescriptionField form={form} />
                  <AmountField form={form} />
                </div>
                <CategorySelector
                  form={form}
                  className="mt-5 flex"
                  iconSize={42}
                  isPredicting={isPredicting}
                  isPredicted={wasPredicted}
                  onCategorySelect={() => {
                    userSelectedCategoryRef.current = true
                    setWasPredicted(false)
                  }}
                />
              </div>
              <div className="mx-auto flex items-center space-x-1">
                <span className="text-md mt-1">Pago por</span>
                <PayerSelector form={form} user={me} partner={partner} />
                <span className="text-md mt-1">dividido</span>
                <SplitRatioSelector form={form} />
              </div>
              <DateSelector form={form} />
              {type === 'update' && <RefundField form={form} />}
            </div>
            <DrawerFooter>
              <div className="mt-4 flex flex-col-reverse gap-2">
                <Button type="submit">Salvar</Button>
                <DrawerClose>
                  <Button type="button" variant="secondary" className="w-full md:w-fit">
                    Cancelar
                  </Button>
                </DrawerClose>
              </div>
            </DrawerFooter>
          </form>
        </Form>
      </DrawerContent>
    </Drawer>
  )
}

export { SaveExpenseButton }

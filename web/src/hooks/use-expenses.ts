/* eslint-disable prettier/prettier */
import { useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { Expense } from '@/domain/expense'
import { privateHttpClient } from '@/http/private-client'
import { ApiError, ApiResponse } from '@/lib/api'

import { useUser } from './use-user'

type GetExpensesApiResponse = ApiResponse<{
  expenses: {
    id: number
    name: string
    amount: number
    refund_amount: number
    description: string
    category_id: number
    payer_id: number
    receiver_id: number
    split_ratio: {
      payer: number
      receiver: number
    }
    split_type: 'equal' | 'proportional' | 'transfer'
    created_at: string
  }[]
  next_token: string
}>

type CreateExpensePayload = Omit<Expense, 'id' | 'description' | 'splitRatio' | 'refundAmount'> & {
  splitType: string
  metadata?: {
    silent?: boolean
  }
}

function useExpenses(search: string) {
  const { user } = useUser()
  const queryClient = useQueryClient()

  const expensesData = useInfiniteQuery<{
    expenses: Expense[]
    nextToken: string
  }>({
    enabled: user?.groupId !== undefined,
    queryKey: ['expenses', search],
    queryFn: async ({ pageParam = '' }): Promise<{ expenses: Expense[]; nextToken: string }> => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetExpensesApiResponse>(
        `/expenses`, {
        params: {
          next_token: pageParam,
          search,
        }
      }
      )

      const expenses = data.expenses?.length
        ? data.expenses?.map<Expense>((expense) => ({
          id: expense.id,
          amount: expense.amount,
          refundAmount: expense.refund_amount,
          categoryId: expense.category_id,
          description: expense.description,
          name: expense.name,
          payerId: expense.payer_id,
          receiverId: expense.receiver_id,
          splitRatio: expense.split_ratio,
          splitType: expense.split_type,
          createdAt: new Date(expense.created_at),
        }))
        : []

      return { expenses, nextToken: data.next_token }
    },
    initialPageParam: '',
    getNextPageParam: (lastPage) => (lastPage?.nextToken !== '' ? lastPage.nextToken : undefined),
  })

  const {
    mutate: createExpense,
    mutateAsync: createExpenseAsync,
    isPending: isCreatingExpense,
  } = useMutation({
    mutationFn: async (payload: CreateExpensePayload) => {
      const { metadata, ...expenseData } = payload

      await privateHttpClient.post(`/expenses`, {
        name: expenseData.name,
        amount: expenseData.amount,
        category_id: expenseData.categoryId,
        split_type: expenseData.splitType,
        payer_id: expenseData.payerId,
        receiver_id: expenseData.receiverId,
        created_at: expenseData.createdAt?.toISOString() ?? new Date().toISOString(),
      })

      return { metadata }
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      queryClient.invalidateQueries({
        queryKey: ['group-balance'],
        exact: true,
      })

      if (variables?.metadata?.silent) {
        return
      }

      toast.success(`Nova despesa criada com sucesso!`)
    },
    onError: (error: ApiError, variables) => {
      let message: string
      switch (error.response?.data?.message) {
        case 'payer income not found':
          message = `O pagador não possui receita cadastrada para dividir proporcionalmente`
          break
        case 'receiver income not found':
          message = `O recebedor não possui receita cadastrada para dividir proporcionalmente`
          break
        default:
          message = error.message
      }

      if (variables?.metadata?.silent) {
        return
      }

      toast.error(`Falha ao criar nova despesa: ${message}`)
    },
  })

  const {
    mutate: updateExpense,
    isPending: isUpdatingExpense,
  } = useMutation({
    mutationFn: async (payload: Partial<Omit<Expense, 'splitRatio'> & { splitType: string }>) => {
      return await privateHttpClient.patch(`/expenses/${payload?.id}`, {
        name: payload.name,
        amount: payload.amount,
        refund_amount: payload.refundAmount,
        category_id: payload.categoryId,
        split_type: payload.splitType,
        payer_id: payload.payerId,
        receiver_id: payload.receiverId,
        created_at: payload.createdAt?.toISOString(),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      queryClient.invalidateQueries({ queryKey: ['expense-details'] })
      queryClient.invalidateQueries({
        queryKey: ['group-balance'],
        exact: true,
      })
      toast.success(`Despesa atualizada com sucesso!`)
    },
    onError: (error: ApiError) => {
      let message: string
      switch (error.response?.data?.message) {
        case 'payer income not found':
          message = `O pagador não possui receita cadastrada para dividir proporcionalmente`
          break
        case 'receiver income not found':
          message = `O recebedor não possui receita cadastrada para dividir proporcionalmente`
          break
        default:
          message = error.message
      }
      toast.error(`Falha ao atualizar despesa: ${message}`)
    },
  })

  const {
    mutate: deleteExpense,
    isPending: isDeletingExpense,
  } = useMutation({
    mutationFn: async (expenseId: number) => await privateHttpClient.delete(`/expenses/${expenseId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      queryClient.invalidateQueries({
        queryKey: ['group-balance'],
        exact: true,
      })
      toast.success(`Despesa deletada com sucesso!`)
    },
    onError: (error) => toast.error(`Falha ao deletar despesa: ${error.message}`),
  })

  return {
    expensesData,
    createExpense,
    createExpenseAsync,
    isCreatingExpense,
    updateExpense,
    isUpdatingExpense,
    deleteExpense,
    isDeletingExpense,
  }
}

export { useExpenses }

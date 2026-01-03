/* eslint-disable prettier/prettier */
import { useQuery } from '@tanstack/react-query'

import { Expense } from '@/domain/expense'
import { privateHttpClient } from '@/http/private-client'
import { ApiResponse } from '@/lib/api'

type GetExpensesApiResponse = ApiResponse<
  {
    id: number
    name: string
    amount: number
    refund_amount: number
    description: string
    category_id: number
    group_id: number
    payer_id: number
    receiver_id: number
    split_ratio: {
      payer: number
      receiver: number
    }
    split_type: 'equal' | 'proportional' | 'transfer'
    created_at: string
    updated_at: string
    deleted_at: string
  }[]
>

function useExpenseDetails(expenseId: number) {
  const { data, isLoading, isError } = useQuery<Expense[]>({
    queryKey: ['expense-details', expenseId],
    staleTime: 0,
    refetchOnWindowFocus: true,
    queryFn: async (): Promise<Expense[]> => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetExpensesApiResponse>(`/expenses/${expenseId}/details`)

      const expenseDetails = data?.length
        ? data?.map<Expense>((expense) => ({
          id: expense.id,
          amount: expense.amount,
          refundAmount: expense.refund_amount,
          categoryId: expense.category_id,
          description: expense.description,
          name: expense.name,
          payerId: expense.payer_id,
          groupId: expense.group_id,
          receiverId: expense.receiver_id,
          splitRatio: expense.split_ratio,
          splitType: expense.split_type,
          createdAt: new Date(expense.created_at),
          updatedAt: new Date(expense.updated_at),
          deletedAt: new Date(expense.deleted_at),
        }))
        : []

      return expenseDetails
    },
  })

  return {
    data,
    isLoading,
    isError,
  }
}

export { useExpenseDetails }

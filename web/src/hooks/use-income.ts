import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { format } from 'date-fns'
import { toast } from 'sonner'

import { Income } from '@/domain/income'
import { privateHttpClient } from '@/http/private-client'
import { ApiResponse } from '@/lib/api'

type GetIncomeApiResponse = ApiResponse<{
  group_id: number
  incomes: {
    id: number
    user_id: number
    type: string
    amount: number
    created_at: string
  }[]
  total: number
  month: number
}>

type CreateIncomePayload = {
  user_id: number
  type: string
  amount: number
  date: Date
}

type UpdateIncomePayload = {
  id: number
  user_id: number
  type?: string
  amount?: number
  date?: Date
}

function useIncome(date?: Date) {
  const queryClient = useQueryClient()

  const { data, isLoading, isError } = useQuery({
    queryKey: ['income', date],
    enabled: !!date,
    queryFn: async () => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetIncomeApiResponse>(`/incomes`, {
        params: { date: format(date ?? new Date(), 'yyyy-MM-dd') },
      })

      const incomes = data.incomes.map<Income>((income) => ({
        id: income.id,
        userId: income.user_id,
        type: income.type as Income['type'],
        amount: income.amount,
        createdAt: new Date(income.created_at),
      }))

      return {
        incomes,
        totalIncome: data.total,
      }
    },
  })

  const { mutate: createIncome } = useMutation({
    mutationFn: async (income: CreateIncomePayload) => {
      await privateHttpClient.post(`/incomes`, {
        user_id: income.user_id,
        type: income.type,
        amount: income.amount,
        created_at: income.date,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['income'] })
      setTimeout(() => queryClient.invalidateQueries({ queryKey: ['expenses'], exact: true }), 2000)
      toast.success('Receita criada com sucesso!')
    },
    onError: (error) => toast.error(`Falha ao criar receita: ${error.message}`),
  })

  const { mutate: updateIncome } = useMutation({
    mutationFn: async (income: UpdateIncomePayload) => {
      await privateHttpClient.patch(`/incomes/${income.id}`, {
        user_id: income.user_id,
        type: income.type,
        amount: income.amount,
        created_at: income.date,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['income'] })
      setTimeout(() => queryClient.invalidateQueries({ queryKey: ['expenses'], exact: true }), 2000)
      toast.success('Receita atualizado com sucesso!')
    },
    onError: (error) => toast.error(`Falha ao atualizar receita: ${error.message}`),
  })

  const { mutate: deleteIncome } = useMutation({
    mutationFn: async (incomeId: number) => {
      await privateHttpClient.delete(`/incomes/${incomeId}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['income'] })
      setTimeout(() => queryClient.invalidateQueries({ queryKey: ['expenses'], exact: true }), 2000)
      toast.success('Receita deletada com sucesso!')
    },
    onError: (error) => toast.error(`Falha ao deletar receita: ${error.message}`),
  })

  return {
    createIncome,
    updateIncome,
    deleteIncome,
    incomes: data?.incomes ?? [],
    totalIncome: data?.totalIncome ?? 0,
    isLoading,
    isError,
  }
}

export { useIncome }

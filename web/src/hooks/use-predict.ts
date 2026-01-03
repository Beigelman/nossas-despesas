import { useMutation } from '@tanstack/react-query'

import { privateHttpClient } from '@/http/private-client'
import { ApiResponse } from '@/lib/api'

type CategoryIDPredictionResponse = ApiResponse<{
  name: string
  amount: number
  category_id: number
}>

type CategoryIDPredictionRequest = {
  name: string
  amount_cents: number
}

function usePredict() {
  const { mutateAsync: predictCategoryID, isPending } = useMutation<
    CategoryIDPredictionResponse,
    unknown,
    CategoryIDPredictionRequest
  >({
    mutationFn: async (req: CategoryIDPredictionRequest) => {
      const { data } = await privateHttpClient.post<CategoryIDPredictionResponse>(`/expenses/predict`, {
        name: req.name,
        amount: req.amount_cents,
      })

      return data
    },
  })

  return {
    predictCategoryID,
    isPredicting: isPending,
  }
}

export { usePredict }

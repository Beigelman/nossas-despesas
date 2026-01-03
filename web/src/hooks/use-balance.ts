import { useQuery } from '@tanstack/react-query'

import { GroupBalance } from '@/domain/group'
import { privateHttpClient } from '@/http/private-client'
import { ApiResponse } from '@/lib/api'

type GetGroupBalanceApiResponse = ApiResponse<{
  group_id: number
  balances: {
    user_id: number
    balance: number
  }[]
}>

function useBalance() {
  const { data, isError, isPending } = useQuery({
    queryKey: ['group-balance'],
    queryFn: async (): Promise<GroupBalance[]> => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetGroupBalanceApiResponse>(`/group/balance`)
      const balances = data.balances?.length
        ? data.balances.map<GroupBalance>((balance) => ({
            balance: balance.balance,
            userId: balance.user_id,
          }))
        : []
      return balances
    },
  })

  return {
    balances: data,
    isLoading: isPending,
    isError,
  }
}

export { useBalance }

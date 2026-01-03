import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useBalance } from '@/hooks/use-balance'
import { useGroup } from '@/hooks/use-group'

import { Skeleton } from '../../../../components/ui/skeleton'
import { UserBalance } from '../../../../components/user-balance'

type OverviewProps = React.ButtonHTMLAttributes<HTMLDivElement>

function Overview({ className }: OverviewProps) {
  const { balances, isLoading } = useBalance()
  const { getMember, me, partner } = useGroup()

  if (!me?.groupId) return null

  if (isLoading) {
    return (
      <div className="flex w-fill-available flex-col gap-2 p-2">
        <Skeleton className="h-4 w-[150px]" />
        {Array.from({ length: 2 }).map((_, i) => (
          <div key={i} className="flex items-center space-x-4">
            <Skeleton className="h-10 w-10 rounded-full" />
            <div className="space-y-2">
              <Skeleton className="h-4 w-[150px]" />
              <Skeleton className="h-4 w-[110px]" />
            </div>
          </div>
        ))}
      </div>
    )
  }

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>Saldos do Grupo</CardTitle>
      </CardHeader>
      <CardContent className="flex items-center justify-between">
        {balances?.map((b) => {
          const member = getMember(b.userId)
          return (
            <UserBalance
              key={b.userId}
              className="mt-2"
              user={{
                name: member?.name ?? '',
                avatar: member?.profileImage ?? '',
                balance: {
                  amount: b.balance,
                  type: b.balance > 0 ? 'income' : 'outcome',
                },
              }}
            />
          )
        })}
        {balances?.length === 0 && partner && (
          <>
            <UserBalance
              className="mt-2"
              user={{
                name: me?.name ?? '',
                avatar: me?.profileImage ?? '',
                balance: {
                  amount: 0,
                  type: 'neutral',
                },
              }}
            />
            <UserBalance
              className="mt-2"
              user={{
                name: partner?.name ?? '',
                avatar: partner?.profileImage ?? '',
                balance: {
                  amount: 0,
                  type: 'neutral',
                },
              }}
            />
          </>
        )}
      </CardContent>
    </Card>
  )
}

export { Overview }

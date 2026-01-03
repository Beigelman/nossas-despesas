import { redirect } from 'next/navigation'

import { getServerSession } from '@/lib/auth/get-server-session'

import { RefreshingSession } from './refreshing-session'

type AcceptGroupInvitePageProps = {
  params: Promise<{ inviteId: string }>
}

export default async function AcceptGroupInvitePage({ params }: AcceptGroupInvitePageProps) {
  const session = await getServerSession()
  const { inviteId } = await params

  const acceptResp = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/group/invite/${inviteId}/accept`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${session?.token}`,
    },
  })
  if (acceptResp.status !== 200) {
    redirect('/')
  }

  return <RefreshingSession refreshToken={session?.refreshToken ?? ''} />
}

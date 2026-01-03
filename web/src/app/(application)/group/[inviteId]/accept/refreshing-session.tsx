'use client'

import { useRouter } from 'next/navigation'
import { signIn } from 'next-auth/react'
import { useEffect } from 'react'
import { toast } from 'sonner'

import { Spinning } from '@/components/ui/spinning'

type RefreshingSessionProps = {
  refreshToken: string
}

function RefreshingSession({ refreshToken }: RefreshingSessionProps) {
  const router = useRouter()

  useEffect(() => {
    signIn('refresh-token', { refreshToken, redirect: false })
      .then(() => {
        toast.success('Convite aceito com sucesso')
        router.replace('/expenses')
      })
      .catch(() => {
        toast.error('Falha ao aceitar convite')
        router.replace('/')
      })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div className="mx-auto flex items-center justify-center">
      <Spinning />
    </div>
  )
}

export { RefreshingSession }

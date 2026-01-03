'use client'

import { Plus, RefreshCw, Send } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'

import { Button } from '@/components/ui/button'
import { Spinning } from '@/components/ui/spinning'
import { useGroup } from '@/hooks/use-group'

import { CreateGroupButton } from '../expenses/components/create-group-button'
import { InviteToGroupButton } from '../expenses/components/invite-to-group-button'

export default function GroupPage() {
  const { me, partner, refresh } = useGroup()
  const router = useRouter()

  useEffect(() => {
    if (me?.groupId && partner) {
      router.replace('/expenses')
    }
  }, [me?.groupId, partner, router])

  if (!me) {
    return (
      <div className="mx-auto flex items-center justify-center">
        <Spinning />
      </div>
    )
  }

  if (!me?.groupId && !partner) {
    return (
      <div className="mx-auto flex items-center justify-center">
        <div className="flex flex-col items-center justify-center">
          <p className="text-center text-xl">Parece que você ainda não faz parte de um grupo</p>
          <p className=" flex flex-wrap items-center justify-center gap-2 text-center">
            Crie um e comece a dividir suas finanças com aquela pessoa especial ♡
          </p>
          <CreateGroupButton>
            <Button className="mt-4 flex items-center justify-center gap-2">
              <Plus /> Criar grupo
            </Button>
          </CreateGroupButton>
        </div>
      </div>
    )
  }

  if (me?.groupId && !partner) {
    return (
      <div className="mx-auto flex items-center justify-center p-3">
        <div className="flex flex-col items-center justify-center">
          <p className="text-center text-xl">Parece que o seu parceiro ainda não entrou na plataforma</p>
          <p className="text-center">Convide-o para que possam começar a sua jornada financeira juntos</p>
          <div className="mt-4 flex items-center gap-1">
            <InviteToGroupButton>
              <Button className="flex items-center justify-center gap-2">
                Convidar <Send size={18} />
              </Button>
            </InviteToGroupButton>
            <Button variant="ghost" size="icon" onClick={() => refresh()}>
              <RefreshCw />
            </Button>
          </div>
        </div>
      </div>
    )
  }
}

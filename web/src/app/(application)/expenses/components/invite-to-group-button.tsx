import { useState } from 'react'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { useGroup } from '@/hooks/use-group'

type InviteToGroupButtonProps = {
  children: React.ReactNode
}

function InviteToGroupButton({ children }: InviteToGroupButtonProps) {
  const [open, setOpen] = useState(false)
  const [email, setEmail] = useState('')
  const { inviteUserToGroup } = useGroup()

  function handleInviteToGroup() {
    inviteUserToGroup(email)
    setOpen(false)
    setEmail('')
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Convidar para o grupo</DialogTitle>
          <DialogDescription>
            Quem você quer convidar para compartilhar as finanças? Mandaremos um e-mail com o convite e assim que
            tivermos um resposta podem começar a sua jornada financeira juntos
          </DialogDescription>
        </DialogHeader>
        <Input value={email} onChange={(e) => setEmail(e.target.value)} />
        <DialogFooter className="mt-4 gap-2">
          <DialogClose>
            <Button variant="secondary" className="w-full md:w-fit">
              Cancelar
            </Button>
          </DialogClose>
          <Button onClick={handleInviteToGroup}>Enviar convite</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export { InviteToGroupButton }

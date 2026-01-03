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
import { useUser } from '@/hooks/use-user'

type CreateGroupButtonProps = {
  children: React.ReactNode
}

function CreateGroupButton({ children }: CreateGroupButtonProps) {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState('')
  const { createGroup } = useUser()

  function handleCreateGroup() {
    createGroup(name)
    setOpen(false)
    setName('')
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Criar novo grupo</DialogTitle>
          <DialogDescription>Que nome esse grupo tão legal vai ter?</DialogDescription>
        </DialogHeader>
        <Input value={name} onChange={(e) => setName(e.target.value)} />
        <DialogFooter className="mt-4 gap-2">
          <DialogClose>
            <Button variant="secondary" className="w-full md:w-fit">
              Cancelar
            </Button>
          </DialogClose>
          <Button onClick={handleCreateGroup}>Criar</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export { CreateGroupButton }

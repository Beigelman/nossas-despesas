'use client'

import { LogOut } from 'lucide-react'
import { signOut } from 'next-auth/react'

import { Button } from './ui/button'

function SignOutButton() {
  return (
    <Button className="flex gap-3" variant={'ghost'} onClick={() => signOut()}>
      <span>Sair</span>
      <LogOut size={20} />
    </Button>
  )
}

export { SignOutButton }

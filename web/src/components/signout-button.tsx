'use client'

import { LogOut } from 'lucide-react'
import { signOut } from 'next-auth/react'
import { useTranslations } from 'next-intl'

import { Button } from './ui/button'

function SignOutButton() {
  const t = useTranslations()

  return (
    <Button className="flex gap-3" variant={'ghost'} onClick={() => signOut()}>
      <span>{t('auth.signOut')}</span>
      <LogOut size={20} />
    </Button>
  )
}

export { SignOutButton }

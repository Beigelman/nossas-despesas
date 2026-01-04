'use client'

import Image from 'next/image'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useTranslations } from 'next-intl'
import { useTheme } from 'next-themes'
import { useMemo } from 'react'

import { cn } from '@/lib/utils'

function MainNav() {
  const pathname = usePathname()
  const { theme } = useTheme()
  const t = useTranslations()

  const systemTheme = useMemo(() => {
    let system = 'light'
    if (typeof window !== 'undefined' && window?.matchMedia) {
      system = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }
    return theme === 'system' ? system : theme
  }, [theme])

  return (
    <div className="mr-4 hidden md:flex">
      <Link href="/" className="mr-6 flex items-center space-x-2">
        <Image
          src={systemTheme === 'dark' ? '/icons/icon-white.png' : '/icons/icon-black.png'}
          alt="logo"
          width={24}
          height={24}
        />
        <span className="hidden sm:inline-block">{t('common.appName')}</span>
      </Link>
      <nav className="flex items-center gap-6 text-sm">
        <Link
          href="/expenses"
          className={cn(
            'transition-colors hover:text-foreground/80',
            pathname === '/expenses' ? 'text-foreground' : 'text-foreground/60',
          )}
        >
          {t('nav.expenses')}
        </Link>
        <Link
          href="/incomes"
          className={cn(
            'transition-colors hover:text-foreground/80',
            pathname?.startsWith('/incomes') ? 'text-foreground' : 'text-foreground/60',
          )}
        >
          {t('nav.incomes')}
        </Link>
        <Link
          href="/insights"
          className={cn(
            'transition-colors hover:text-foreground/80',
            pathname?.startsWith('/insights') ? 'text-foreground' : 'text-foreground/60',
          )}
        >
          {t('nav.insights')}
        </Link>
      </nav>
    </div>
  )
}

export { MainNav }

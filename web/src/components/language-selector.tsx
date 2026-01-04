'use client'

import { Globe } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useLocale } from 'next-intl'
import { useState } from 'react'

import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { type Locale, locales } from '@/i18n/config'

const localeNames: Record<Locale, string> = {
  'pt-BR': 'Português',
  en: 'English',
}

function LanguageSelector() {
  const locale = useLocale() as Locale
  const router = useRouter()
  const [isChanging, setIsChanging] = useState(false)

  function changeLocale(newLocale: Locale) {
    if (newLocale === locale || isChanging) return

    setIsChanging(true)
    // Set cookie
    document.cookie = `NEXT_LOCALE=${newLocale}; path=/; max-age=${60 * 60 * 24 * 365}`
    // Reload page to apply new locale
    router.refresh()
    setTimeout(() => setIsChanging(false), 1000)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="h-9 w-9">
          <Globe className="h-4 w-4" />
          <span className="sr-only">Change language</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {locales.map((loc) => (
          <DropdownMenuItem
            key={loc}
            onClick={() => changeLocale(loc)}
            disabled={isChanging || loc === locale}
            className={loc === locale ? 'bg-accent' : ''}
          >
            {localeNames[loc]}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

export { LanguageSelector }

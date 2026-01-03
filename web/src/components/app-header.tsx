'use client'

import { usePathname } from 'next/navigation'

import { MainNav } from '@/components/main-nav'
import { MobileNav } from '@/components/mobile-nav'
import { ModeToggle } from '@/components/mode-toggle'

import { SearchBar } from './search-bar'
import { UserNav } from './user-nav'

function AppHeader() {
  const pathname = usePathname()

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 max-w-screen-2xl items-center">
        <MainNav />
        <MobileNav />
        <div className="flex flex-1 items-center justify-end space-x-2 md:justify-end">
          {pathname?.startsWith('/expenses') && <SearchBar />}
          <ModeToggle />
          <UserNav />
        </div>
      </div>
    </header>
  )
}

export { AppHeader }

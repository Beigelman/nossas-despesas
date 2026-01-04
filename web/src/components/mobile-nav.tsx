'use client'

import { Menu } from 'lucide-react'
import Image from 'next/image'
import Link, { LinkProps } from 'next/link'
import { useRouter } from 'next/navigation'
import { useTranslations } from 'next-intl'
import { useTheme } from 'next-themes'
import * as React from 'react'

import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { cn } from '@/lib/utils'

function MobileNav() {
  const [open, setOpen] = React.useState(false)
  const { theme } = useTheme()
  const t = useTranslations()

  const systemTheme = React.useMemo(() => {
    let system = 'light'
    if (typeof window !== 'undefined' && window?.matchMedia) {
      system = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }
    return theme === 'system' ? system : theme
  }, [theme])

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          variant="ghost"
          className="mr-2 px-0 text-base hover:bg-transparent focus-visible:bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0 md:hidden"
        >
          <Menu />
          <span className="sr-only">Toggle Menu</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="pr-0">
        <MobileLink href="/" className="flex items-center" onOpenChange={setOpen}>
          <Image
            src={systemTheme === 'dark' ? '/icons/icon-white.png' : '/icons/icon-black.png'}
            alt="logo"
            width={24}
            height={24}
          />
          <span className="ml-2 font-bold">{t('common.appName')}</span>
        </MobileLink>
        <ScrollArea className="my-4 h-[calc(100vh-8rem)] pb-10 pl-6">
          <div className="flex flex-col space-y-3">
            <MobileLink href={'/expenses'} onOpenChange={setOpen}>
              {t('nav.expenses')}
            </MobileLink>
            <MobileLink href={'/incomes'} onOpenChange={setOpen}>
              {t('nav.incomes')}
            </MobileLink>
            <MobileLink href={'/insights'} onOpenChange={setOpen}>
              {t('nav.insights')}
            </MobileLink>
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  )
}

interface MobileLinkProps extends LinkProps {
  onOpenChange?: (open: boolean) => void
  children: React.ReactNode
  className?: string
}

function MobileLink({ href, onOpenChange, className, children, ...props }: MobileLinkProps) {
  const router = useRouter()
  return (
    <Link
      href={href}
      onClick={() => {
        router.push(href.toString())
        onOpenChange?.(false)
      }}
      className={cn(className)}
      {...props}
    >
      {children}
    </Link>
  )
}

export { MobileNav }

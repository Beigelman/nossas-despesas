'use client'

import { LogIn } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'
import { useTranslations } from 'next-intl'
import { useTheme } from 'next-themes'
import { useMemo } from 'react'

import { Button } from '@/components/ui/button'

export default function LandingPage() {
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
    <div className="flex min-h-screen flex-col">
      <header className="flex h-14 items-center px-4 lg:px-6">
        <Link className="flex items-center justify-center" href="#">
          <Image
            src={systemTheme === 'dark' ? '/icons/icon-white.png' : '/icons/icon-black.png'}
            alt="logo"
            width={24}
            height={24}
          />
          <span className="ml-2">{t('common.appName')}</span>
        </Link>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <Button variant="ghost" className="flex gap-3">
            <Link className="text-sm font-medium underline-offset-4 hover:underline" href="/login">
              {t('landing.signUpLogIn')}
            </Link>
            <LogIn className="h-6 w-6" />
          </Button>
        </nav>
      </header>
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center space-y-4 text-center">
              <div className="space-y-2">
                <h1 className="lg:text-6xl/none text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
                  {t('landing.title')}
                </h1>
                <p className="mx-auto max-w-[700px] text-gray-500 dark:text-gray-400 md:text-xl">
                  {t('landing.subtitle')}
                </p>
              </div>
              <div className="space-x-4">
                <Link
                  className="inline-flex h-9 items-center justify-center rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-gray-50 shadow transition-colors hover:bg-gray-900/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-gray-950 disabled:pointer-events-none disabled:opacity-50 dark:bg-gray-50 dark:text-gray-900 dark:hover:bg-gray-50/90 dark:focus-visible:ring-gray-300"
                  href="/login"
                >
                  {t('landing.getStarted')}
                </Link>
              </div>
            </div>
          </div>
        </section>
        <section className="w-full bg-gray-100 py-12 dark:bg-gray-800 md:py-24 lg:py-32">
          <div className="container px-4 md:px-6">
            <div className="grid items-center gap-6 lg:grid-cols-3 lg:gap-12 xl:grid-cols-3">
              <div className="flex flex-col justify-center space-y-4">
                <div className="space-y-2">
                  <div className="inline-block rounded-lg bg-gray-100 px-3 py-1 text-sm dark:bg-gray-800">
                    {t('landing.features.easyToUse.badge')}
                  </div>
                  <h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">
                    {t('landing.features.easyToUse.title')}
                  </h2>
                  <p className="max-w-[600px] text-gray-500 dark:text-gray-400 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                    {t('landing.features.easyToUse.description')}
                  </p>
                </div>
              </div>
              <div className="flex flex-col justify-center space-y-4">
                <div className="space-y-2">
                  <div className="inline-block rounded-lg bg-gray-100 px-3 py-1 text-sm dark:bg-gray-800">
                    {t('landing.features.fairAndTransparent.badge')}
                  </div>
                  <h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">
                    {t('landing.features.fairAndTransparent.title')}
                  </h2>
                  <p className="max-w-[600px] text-gray-500 dark:text-gray-400 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                    {t('landing.features.fairAndTransparent.description')}
                  </p>
                </div>
              </div>
              <div className="flex flex-col justify-center space-y-4">
                <div className="space-y-2">
                  <div className="inline-block rounded-lg bg-gray-100 px-3 py-1 text-sm dark:bg-gray-800">
                    {t('landing.features.secureAndReliable.badge')}
                  </div>
                  <h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">
                    {t('landing.features.secureAndReliable.title')}
                  </h2>
                  <p className="max-w-[600px] text-gray-500 dark:text-gray-400 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                    {t('landing.features.secureAndReliable.description')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>
      <footer className="flex w-full shrink-0 flex-col items-center gap-2 border-t px-4 py-6 sm:flex-row md:px-6">
        <p className="text-xs text-gray-500 dark:text-gray-400">{t('landing.footer.copyright')}</p>
        <nav className="flex gap-4 sm:ml-auto sm:gap-6">
          <Link className="text-xs underline-offset-4 hover:underline" href="#">
            {t('landing.footer.termsOfService')}
          </Link>
          <Link className="text-xs underline-offset-4 hover:underline" href="#">
            {t('landing.footer.privacy')}
          </Link>
        </nav>
      </footer>
    </div>
  )
}

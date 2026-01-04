import { cookies } from 'next/headers'
import { getRequestConfig } from 'next-intl/server'

import { defaultLocale, type Locale } from './config'

export default getRequestConfig(async ({ requestLocale }) => {
  // Get locale from request or cookie, or use default
  let locale = (await requestLocale) as Locale | undefined

  if (!locale) {
    const cookieStore = await cookies()
    locale = (cookieStore.get('NEXT_LOCALE')?.value || defaultLocale) as Locale
  }

  return {
    locale,
    messages: (await import(`@/messages/${locale}.json`)).default,
  }
})

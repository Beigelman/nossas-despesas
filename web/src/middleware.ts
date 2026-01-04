import type { NextRequest } from 'next/server'
import { NextResponse } from 'next/server'
import createMiddleware from 'next-intl/middleware'

import { defaultLocale, locales } from '@/i18n/config'

const intlMiddleware = createMiddleware({
  locales,
  defaultLocale,
  localePrefix: 'never', // No prefix in URL
  localeDetection: true, // Detect from Accept-Language header
})

export function middleware(request: NextRequest) {
  // Apply intl middleware first
  const intlResponse = intlMiddleware(request)

  // Add custom header for URL
  const requestHeaders = new Headers(request.headers)
  requestHeaders.set('x-url', request.url)

  // Create new response - always use NextResponse.next to avoid rewrite issues
  // Since we use localePrefix: 'never', we don't want the middleware to rewrite URLs
  const response = NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  })

  // Copy cookies from intl middleware
  intlResponse.cookies.getAll().forEach((cookie) => {
    response.cookies.set(cookie.name, cookie.value, cookie)
  })

  // Copy locale header from intl middleware
  const localeHeader = intlResponse.headers.get('x-next-intl-locale')
  if (localeHeader) {
    response.headers.set('x-next-intl-locale', localeHeader)
  }

  return response
}

export const config = {
  matcher: ['/((?!api|_next|_vercel|.*\\..*).*)'],
}

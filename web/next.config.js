/* eslint-disable @typescript-eslint/no-var-requires */
/** @type {import('next').NextConfig} */
const isPWAEnabled = process.env.NEXT_DISABLE_PWA !== 'true'
const withPWA = require('@ducanh2912/next-pwa').default({
  dest: 'public',
  disable: process.env.NODE_ENV === 'development' || !isPWAEnabled,
  register: true,
  scope: '/app',
})

const createNextIntlPlugin = require('next-intl/plugin')('src/i18n/request.ts')

const nextConfig = {
  // next.js config
}

const configWithPWA = isPWAEnabled ? withPWA(nextConfig) : nextConfig

module.exports = createNextIntlPlugin(configWithPWA)

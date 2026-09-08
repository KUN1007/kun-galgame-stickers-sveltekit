import { existsSync, readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import tailwindcss from '@tailwindcss/vite'

const unquote = (value: string) => {
  if (
    (value.startsWith('"') && value.endsWith('"')) ||
    (value.startsWith("'") && value.endsWith("'"))
  ) {
    return value.slice(1, -1)
  }
  return value
}

const applyRootEnv = (raw: string) => {
  const parsed: Record<string, string> = {}
  for (const line of raw.split('\n')) {
    const trimmed = line.trim()
    if (!trimmed || trimmed.startsWith('#')) continue
    const eq = trimmed.indexOf('=')
    if (eq < 1) continue
    parsed[trimmed.slice(0, eq).trim()] = unquote(trimmed.slice(eq + 1).trim())
  }
  const lookup = (name: string) => parsed[name] ?? process.env[name] ?? ''
  for (const [key, value] of Object.entries(parsed)) {
    parsed[key] = value.replace(/\$\{([A-Z0-9_]+)\}/g, (_, name: string) => lookup(name))
  }
  for (const [key, value] of Object.entries(parsed)) {
    const current = process.env[key]
    if (current === undefined || /\$\{[A-Z0-9_]+\}/.test(current)) {
      process.env[key] = value
    }
  }
}

const rootEnv = fileURLToPath(new URL('../../.env', import.meta.url))
if (existsSync(rootEnv)) applyRootEnv(readFileSync(rootEnv, 'utf8'))

const envValue = (...keys: string[]) => {
  for (const key of keys) {
    const value = process.env[key]
    if (value && !/\$\{[A-Z0-9_]+\}/.test(value)) return unquote(value)
  }
  return ''
}

export default defineNuxtConfig({
  extends: ['@kungal/ui-nuxt'],

  compatibilityDate: '2026-08-08',
  devtools: { enabled: false },

  modules: ['@nuxt/eslint', '@nuxtjs/color-mode', '@nuxtjs/i18n', '@nuxtjs/sitemap'],

  css: ['~/assets/css/main.css'],

  vite: {
    plugins: [tailwindcss()]
  },

  imports: {
    dirs: ['~/features/**'],
    presets: [{ from: '@kungal/ui-core', imports: ['cn'] }]
  },

  devServer: {
    host: '127.0.0.1',
    port: 5173
  },

  runtimeConfig: {
    apiBaseUrl: envValue('NUXT_API_BASE_URL') || 'http://127.0.0.1:9421',
    public: {
      apiBaseUrl: envValue('NUXT_PUBLIC_API_BASE_URL') || 'http://127.0.0.1:9421',
      siteUrl: envValue('NUXT_PUBLIC_SITE_URL') || 'https://sticker.kungal.com',
      oauthServerUrl:
        envValue('NUXT_PUBLIC_OAUTH_SERVER_URL', 'KUN_OAUTH_SERVER_URL') ||
        'http://127.0.0.1:9277/api/v1',
      oauthFrontendUrl:
        envValue('NUXT_PUBLIC_OAUTH_FRONTEND_URL', 'KUN_OAUTH_WEB_URL') ||
        'http://127.0.0.1:9420',
      oauthClientId:
        envValue('NUXT_PUBLIC_OAUTH_CLIENT_ID', 'KUN_OAUTH_CLIENT_ID') ||
        'c5cd7b074804ba134934eb6c175a8f4d',
      oauthRedirectUri:
        envValue('NUXT_PUBLIC_OAUTH_REDIRECT_URI', 'KUN_OAUTH_REDIRECT_URI') ||
        'http://127.0.0.1:5173/auth/callback'
    }
  },

  i18n: {
    strategy: 'prefix_except_default',
    defaultLocale: 'zh-cn',
    langDir: 'locales',
    baseUrl: process.env.NUXT_PUBLIC_SITE_URL || 'https://sticker.kungal.com',
    locales: [
      { code: 'zh-cn', language: 'zh-CN', name: '中文', file: 'zh-cn.json' },
      { code: 'en', language: 'en-US', name: 'English', file: 'en.json' },
      { code: 'ja', language: 'ja-JP', name: '日本語', file: 'ja.json' }
    ],
    detectBrowserLanguage: false
  },

  colorMode: {
    preference: 'system',
    fallback: 'light',
    classPrefix: 'kun-',
    classSuffix: '-mode',
    storageKey: 'kun-sticker-theme'
  },

  site: {
    url: process.env.NUXT_PUBLIC_SITE_URL || 'https://sticker.kungal.com'
  },

  sitemap: {
    exclude: ['/auth/**'],
    sources: ['/api/__sitemap__/urls'],
    excludeAppSources: true,
    cacheMaxAgeSeconds: 60 * 60 * 6,
    defaults: { changefreq: 'weekly', priority: 0.7 }
  },

  app: {
    head: {
      link: [{ rel: 'icon', type: 'image/webp', href: '/favicon.webp' }]
    }
  },

  eslint: {
    config: { stylistic: false }
  }
})

import i18n from 'sveltekit-i18n'
import type { Config } from 'sveltekit-i18n'

export const config: Config = {
  loaders: [
    {
      locale: 'en',
      key: 'about',
      loader: async () => (await import('./en/about.json')).default
    },
    {
      locale: 'en',
      key: 'game',
      loader: async () => (await import('./en/game.json')).default
    },
    {
      locale: 'en',
      key: 'header',
      loader: async () => (await import('./en/header.json')).default
    },
    {
      locale: 'en',
      key: 'home',
      loader: async () => (await import('./en/home.json')).default
    },
    {
      locale: 'en',
      key: 'lass',
      loader: async () => (await import('./en/lass.json')).default
    },
    {
      locale: 'en',
      key: 'seo',
      loader: async () => (await import('./en/seo.json')).default
    },
    {
      locale: 'en',
      key: 'sticker',
      loader: async () => (await import('./en/sticker.json')).default
    },
    {
      locale: 'zh',
      key: 'about',
      loader: async () => (await import('./zh/about.json')).default
    },
    {
      locale: 'zh',
      key: 'game',
      loader: async () => (await import('./zh/game.json')).default
    },
    {
      locale: 'zh',
      key: 'header',
      loader: async () => (await import('./zh/header.json')).default
    },
    {
      locale: 'zh',
      key: 'home',
      loader: async () => (await import('./zh/home.json')).default
    },
    {
      locale: 'zh',
      key: 'lass',
      loader: async () => (await import('./zh/lass.json')).default
    },
    {
      locale: 'zh',
      key: 'seo',
      loader: async () => (await import('./zh/seo.json')).default
    },
    {
      locale: 'zh',
      key: 'sticker',
      loader: async () => (await import('./zh/sticker.json')).default
    }
  ]
}

export const { t, loading, locales, locale, loadTranslations } = new i18n(config)

loading.subscribe(($loading) => $loading)

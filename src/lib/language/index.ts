import i18n from 'sveltekit-i18n'
import type { Config } from 'sveltekit-i18n'

export const config: Config = {
  loaders: [
    {
      locale: 'en-us',
      key: 'about',
      loader: async () => (await import('./en-us/about.json')).default
    },
    {
      locale: 'en-us',
      key: 'header',
      loader: async () => (await import('./en-us/header.json')).default
    },
    {
      locale: 'en-us',
      key: 'home',
      loader: async () => (await import('./en-us/home.json')).default
    },
    {
      locale: 'en-us',
      key: 'seo',
      loader: async () => (await import('./en-us/seo.json')).default
    },
    {
      locale: 'en-us',
      key: 'sticker',
      loader: async () => (await import('./en-us/sticker.json')).default
    },

    /* ---------------- Chinese --------------- */

    {
      locale: 'zh-cn',
      key: 'about',
      loader: async () => (await import('./zh-cn/about.json')).default
    },
    {
      locale: 'zh-cn',
      key: 'header',
      loader: async () => (await import('./zh-cn/header.json')).default
    },
    {
      locale: 'zh-cn',
      key: 'home',
      loader: async () => (await import('./zh-cn/home.json')).default
    },
    {
      locale: 'zh-cn',
      key: 'seo',
      loader: async () => (await import('./zh-cn/seo.json')).default
    },
    {
      locale: 'zh-cn',
      key: 'sticker',
      loader: async () => (await import('./zh-cn/sticker.json')).default
    }
  ]
}

export const { t, loading, locales, locale, loadTranslations } = new i18n(config)

loading.subscribe(($loading) => $loading)

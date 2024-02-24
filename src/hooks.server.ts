import { PUBLIC_KUN_STICKER_THEME, PUBLIC_KUN_LANGUAGE } from '$env/static/public'
import { locales } from '$lib/language'
import type { Handle } from '@sveltejs/kit'

export const handle: Handle = async ({ event, resolve }) => {
  const { url, locals, cookies, request } = event
  const { pathname } = url

  const lang = request.headers.get('accept-language') || 'en'
  const cookieLanguage = (cookies.get(PUBLIC_KUN_LANGUAGE) as App.KunLanguage) || 'en'
  const kunLanguage = () => {
    if (lang === 'zh') {
      return 'zh'
    }
    return 'en'
  }
  locals.language = cookieLanguage ? cookieLanguage : kunLanguage()
  cookies.set(PUBLIC_KUN_LANGUAGE, locals.language, {
    path: '/',
    maxAge: 604800,
    httpOnly: false
  })

  locals.theme = (cookies.get(PUBLIC_KUN_STICKER_THEME) as App.KunTheme) || 'system'
  cookies.set(PUBLIC_KUN_STICKER_THEME, locals.theme, {
    path: '/',
    maxAge: 604800,
    httpOnly: false
  })

  const supportedLocales = locales.get().map((l) => l.toLowerCase())
  let locale = supportedLocales.find(
    (l) => l === `${pathname.match(/[^/]+?(?=\/|$)/)}`.toLowerCase()
  )

  if (!locale) {
    locale = kunLanguage()
    return new Response(undefined, { headers: { location: `/${locale}${pathname}` }, status: 301 })
  }

  return await resolve(event, {
    transformPageChunk: ({ html }) => {
      const replacedHtml = html
        .replace('%language%', event.locals.language)
        .replace('%kun-sticker-theme%', event.locals.theme)
      return replacedHtml
    }
  })
}

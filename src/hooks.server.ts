import { PUBLIC_KUN_STICKER_THEME, PUBLIC_KUN_LANGUAGE } from '$env/static/public'
import type { Handle } from '@sveltejs/kit'

export const handle: Handle = async ({ event, resolve }) => {
  const { locals, cookies, request } = event

  const lang = request.headers.get('accept-language') || 'en'
  const kunLanguage = () => {
    if (lang === 'zh') {
      return 'zh'
    }
    return 'en'
  }
  locals.language = kunLanguage()
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

  return await resolve(event, {
    transformPageChunk: ({ html }) => {
      const replacedHtml = html
        .replace('%language%', event.locals.language)
        .replace('%kun-sticker-theme%', event.locals.theme)
      return replacedHtml
    }
  })
}

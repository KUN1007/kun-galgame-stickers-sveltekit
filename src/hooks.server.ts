import { PUBLIC_KUN_STICKER_THEME } from '$env/static/public'
import type { Handle } from '@sveltejs/kit'

export const handle: Handle = async ({ event, resolve }) => {
  const { locals, cookies } = event

  locals.theme = (cookies.get(PUBLIC_KUN_STICKER_THEME) as App.KunTheme) || 'system'
  cookies.set(PUBLIC_KUN_STICKER_THEME, locals.theme, {
    path: '/',
    maxAge: 604800,
    httpOnly: false
  })

  return await resolve(event, {
    transformPageChunk: ({ html }) => {
      const replacedHtml = html.replace('%kun-sticker-theme%', event.locals.theme)
      return replacedHtml
    }
  })
}

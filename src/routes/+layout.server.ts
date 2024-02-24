import { loadTranslations, translations } from '$lib/language'

export const load = async ({ url, locals }) => {
  const { pathname } = url
  const { language } = locals

  const route = pathname.replace(new RegExp(`^/${language}`), '')

  await loadTranslations(language, route)

  return { i18n: { route, language }, translations: translations.get() }
}

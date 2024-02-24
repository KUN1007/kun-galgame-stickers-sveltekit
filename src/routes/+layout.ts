import { addTranslations, setLocale, setRoute } from '$lib/language'

export const load = async ({ data }) => {
  const { i18n, translations } = data
  const { language, route } = i18n

  addTranslations(translations)

  await setRoute(route)
  await setLocale(language)

  return i18n
}

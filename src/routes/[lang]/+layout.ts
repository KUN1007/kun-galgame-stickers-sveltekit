import { loadTranslations } from '~/lib/language'

interface KunData {
  theme: App.KunTheme
  language: App.KunLanguage
}

export const load = async ({ data }: { data: KunData }) => {
  await loadTranslations(data.language)
  return { theme: data.theme, language: data.language }
}

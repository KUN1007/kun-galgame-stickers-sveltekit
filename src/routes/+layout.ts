import { loadTranslations } from '~/lib/language'
import type { Load } from '@sveltejs/kit'

interface KunData {
  theme: App.KunTheme
  language: App.KunLanguage
}

export const load: Load<KunData> = async ({ data }) => {
  await loadTranslations(data?.language)
  return { theme: data?.theme, language: data?.language }
}

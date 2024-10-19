import { getContext, setContext } from 'svelte'
import { writable } from 'svelte/store'
import { locale } from '~/lib/language'
import { PUBLIC_KUN_LANGUAGE } from '$env/static/public'

export const createLanguageStore = (language: Language) => {
  const store = writable(language)

  return {
    subscribe: store.subscribe,
    change(lang: Language) {
      locale.set(lang)
      document.cookie = `${PUBLIC_KUN_LANGUAGE}=${lang}; path=/; Max-Age=604800`
      store.set(lang)
    }
  }
}

export function setLanguageContext(language: Language) {
  return setContext('language', createLanguageStore(language))
}

export function getLanguageContext() {
  return getContext<ReturnType<typeof setLanguageContext>>('language')
}

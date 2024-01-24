import { getContext, setContext } from 'svelte'
import { writable } from 'svelte/store'
import { locale } from '~/lib/language'
import { PUBLIC_KUN_LANGUAGE } from '$env/static/public'

export const createLanguageStore = (language: App.KunLanguage) => {
  const store = writable(language)

  return {
    subscribe: store.subscribe,
    change(lang: App.KunLanguage) {
      locale.set(lang)
      document.cookie = `${PUBLIC_KUN_LANGUAGE}=${lang}; path=/; Max-Age=604800`
      store.set(lang)
    }
  }
}

export function setLanguageContext(language: App.KunLanguage) {
  return setContext('language', createLanguageStore(language))
}

export function getLanguageContext() {
  return getContext<ReturnType<typeof setLanguageContext>>('language')
}

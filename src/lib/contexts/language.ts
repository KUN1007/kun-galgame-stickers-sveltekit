import { getContext, setContext } from 'svelte'
import { writable } from 'svelte/store'

export function setLanguageContext(initial: App.KunLanguage) {
  const lang = writable(initial)
  return setContext('language', {
    lang
  })
}

export function getLanguageContext() {
  return getContext<ReturnType<typeof setLanguageContext>>('language')
}

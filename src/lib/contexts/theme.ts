import { getContext, setContext } from 'svelte'
import { derived, writable } from 'svelte/store'

import { PUBLIC_KUN_STICKER_THEME } from '$env/static/public'

const getPrefersColorTheme = () => {
  return window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

const createColorThemeStore = (theme: App.KunTheme) => {
  const store = writable(theme)

  const preferred = derived(store, (theme) => (theme === 'system' ? getPrefersColorTheme() : theme))

  return {
    subscribe: store.subscribe,
    change(scheme: App.KunTheme) {
      document.documentElement.dataset.colorScheme = scheme
      document.cookie = `${PUBLIC_KUN_STICKER_THEME}=${scheme}; path=/; Max-Age=604800`
      store.set(scheme)
    },
    preferred
  }
}

export const setColorSchemeContext = (initial: App.KunTheme) => {
  return setContext('theme', createColorThemeStore(initial))
}

export const getColorSchemeContext = () => {
  return getContext<ReturnType<typeof setColorSchemeContext>>('theme')
}

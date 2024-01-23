// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
  namespace App {
    declare type KunTheme = 'light' | 'dark' | 'system'
    declare type KunLanguage = 'en' | 'zh'

    // interface Error {}
    interface Locals {
      theme: KunTheme
      language: KunLanguage
    }
    // interface PageData {}
    // interface PageState {}
    // interface Platform {}
  }
}

export {}

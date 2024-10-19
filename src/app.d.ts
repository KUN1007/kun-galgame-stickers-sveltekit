// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
  namespace App {
    declare type KunTheme = 'light' | 'dark' | 'system'

    // interface Error {}
    interface Locals {
      theme: KunTheme
      language: Language
    }
    // interface PageData {}
    // interface PageState {}
    // interface Platform {}
  }
}

export {}

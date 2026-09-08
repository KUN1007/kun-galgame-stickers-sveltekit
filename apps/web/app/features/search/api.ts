import type { CatalogCharacter, Pack } from '~/features/pack/types'

export interface SearchResults {
  packs: Pack[]
  characters: CatalogCharacter[]
}

export const quickSearch = async (q: string): Promise<SearchResults> => {
  const trimmed = q.trim()
  if (!trimmed) return { packs: [], characters: [] }
  const data = await kunFetchOrNull<SearchResults>(`/search?q=${encodeURIComponent(trimmed)}`)
  return data ?? { packs: [], characters: [] }
}

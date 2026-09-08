import type { CatalogCharacter, CatalogWork, CharacterPage } from '~/features/pack/types'

/**
 * Catalog reads go through this site's own API, never straight to infra: the
 * application key is a server secret. The two picker calls need a session --
 * only an author composing a pack has any use for them.
 */
export const searchCatalogWorks = async (q: string): Promise<CatalogWork[]> => {
  const trimmed = q.trim()
  if (!trimmed) return []
  const data = await kunFetchOrNull<{ works: CatalogWork[] }>(
    `/catalog/works?q=${encodeURIComponent(trimmed)}`
  )
  return data?.works ?? []
}

export const fetchWorkRoster = async (workId: number): Promise<CatalogCharacter[]> => {
  const data = await kunFetchOrNull<{ characters: CatalogCharacter[] }>(
    `/catalog/works/${workId}/characters`
  )
  return data?.characters ?? []
}

export const fetchCharacter = (characterId: number | string): Promise<CharacterPage | null> =>
  kunFetchOrNull<CharacterPage>(`/characters/${characterId}`)

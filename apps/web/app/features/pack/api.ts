import type { Pack, PackDetail, PackListPage, MultilingualText } from './types'

export interface PackQuery {
  page?: number
  limit?: number
  sort?: 'new' | 'hot'
  q?: string
  tag?: string
  rating?: 'all' | 'sfw'
  official?: boolean
  /** Keep only packs that declare a catalog game. */
  linked?: boolean
  /** List the packs made from one game. */
  work?: number
}

const toParams = (query: PackQuery): Record<string, string> => {
  const params: Record<string, string> = {}
  if (query.page && query.page > 1) params.page = String(query.page)
  if (query.limit) params.limit = String(query.limit)
  if (query.sort) params.sort = query.sort
  if (query.q) params.q = query.q
  if (query.tag) params.tag = query.tag
  if (query.rating) params.rating = query.rating
  if (query.official) params.official = '1'
  if (query.linked) params.linked = '1'
  if (query.work) params.work = String(query.work)
  return params
}

const withQuery = (path: string, query: PackQuery): string => {
  const search = new URLSearchParams(toParams(query)).toString()
  return search ? `${path}?${search}` : path
}

export const fetchPacks = (query: PackQuery = {}): Promise<PackListPage | null> =>
  kunFetchOrNull<PackListPage>(withQuery('/packs', query))

export const fetchUserPacks = (uid: number, query: PackQuery = {}): Promise<PackListPage | null> =>
  kunFetchOrNull<PackListPage>(withQuery(`/users/${uid}/packs`, query))

export const fetchPack = (packId: string): Promise<PackDetail | null> =>
  kunFetchOrNull<PackDetail>(`/packs/${packId}`)

export const fetchMyPacks = (query: PackQuery = {}): Promise<PackListPage | null> =>
  kunFetchOrNull<PackListPage>(withQuery('/me/packs', query))

export interface CreatePackBody {
  title: MultilingualText
  description?: MultilingualText
  content_rating?: number
  tags?: string[]
  catalog_work_id?: number
}

export const createPack = (body: CreatePackBody): Promise<Pack> =>
  kunFetch<Pack>('/me/packs', { method: 'POST', body })

export interface PatchPackBody {
  title?: MultilingualText
  description?: MultilingualText
  content_rating?: number
  cover_sticker_id?: string
  tags?: string[]
  /** 0 clears the link -- an omitted field means "leave it alone". */
  catalog_work_id?: number
}

export const patchPack = (packId: string, body: PatchPackBody): Promise<Pack> =>
  kunFetch<Pack>(`/me/packs/${packId}`, { method: 'PATCH', body })

export const deletePack = (packId: string): Promise<unknown> =>
  kunFetch(`/me/packs/${packId}`, { method: 'DELETE' })

export const publishPack = (packId: string): Promise<Pack> =>
  kunFetch<Pack>(`/me/packs/${packId}/publish`, { method: 'POST' })

export const unpublishPack = (packId: string): Promise<Pack> =>
  kunFetch<Pack>(`/me/packs/${packId}/unpublish`, { method: 'POST' })

export const packDownloadUrl = (packId: string): string =>
  `${useRuntimeConfig().public.apiBaseUrl}/api/v1/packs/${packId}/download`

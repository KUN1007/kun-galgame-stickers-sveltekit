import type { ImageUpload, MultilingualText, Sticker } from '~/features/pack/types'

export const uploadPackImage = (packId: string, file: File): Promise<ImageUpload> => {
  const body = new FormData()
  body.append('file', file)
  return kunFetch<ImageUpload>(`/me/packs/${packId}/images`, { method: 'POST', body })
}

export interface AddStickerBody {
  image_hash: string
  width?: number
  height?: number
  game?: MultilingualText
  character_name?: MultilingualText
  vndb_id?: number
  note?: string
  catalog_work_id?: number
  catalog_character_id?: number
}

export const addSticker = (packId: string, body: AddStickerBody): Promise<Sticker> =>
  kunFetch<Sticker>(`/me/packs/${packId}/stickers`, { method: 'POST', body })

export interface PatchStickerBody {
  game?: MultilingualText
  character_name?: MultilingualText
  vndb_id?: number
  note?: string
  /** 0 clears the link -- an omitted field means "leave it alone". */
  catalog_work_id?: number
  catalog_character_id?: number
}

export const patchSticker = (
  packId: string,
  stickerId: string,
  body: PatchStickerBody
): Promise<Sticker> =>
  kunFetch<Sticker>(`/me/packs/${packId}/stickers/${stickerId}`, { method: 'PATCH', body })

export const deleteSticker = (packId: string, stickerId: string): Promise<unknown> =>
  kunFetch(`/me/packs/${packId}/stickers/${stickerId}`, { method: 'DELETE' })

export const reorderStickers = (packId: string, stickerIds: string[]): Promise<unknown> =>
  kunFetch(`/me/packs/${packId}/stickers/order`, {
    method: 'PUT',
    body: { sticker_ids: stickerIds }
  })

export const fetchSticker = (stickerId: string): Promise<Sticker | null> =>
  kunFetchOrNull<Sticker>(`/stickers/${stickerId}`)

export const stickerDownloadUrl = (stickerId: string): string =>
  `${useRuntimeConfig().public.apiBaseUrl}/api/v1/stickers/${stickerId}/download`

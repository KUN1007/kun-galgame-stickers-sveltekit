export type MultilingualText = Record<string, string | undefined>

export interface Author {
  id: number
  name: string
  avatar: string
}

export interface Tag {
  id: string
  slug: string
  name: MultilingualText
  pack_count: number
}

export interface Sticker {
  id: string
  pack_id: string
  position: number
  width: number
  height: number
  game: MultilingualText
  character_name: MultilingualText
  vndb_id?: number
  note: string
  image_url: string
  thumb_url: string
}

export interface Pack {
  id: string
  status: number
  is_official: boolean
  content_rating: number
  title: MultilingualText
  description: MultilingualText
  cover_url: string
  cover_thumb_url: string
  sticker_count: number
  view_count: number
  download_count: number
  author: Author
  tags: Tag[]
  created_at: string
  updated_at: string
  published_at?: string
}

export interface PackDetail extends Pack {
  stickers: Sticker[]
}

export interface PackListPage {
  packs: Pack[]
  total: number
}

export interface ImageUpload {
  hash: string
  image_url: string
  thumb_url: string
  width: number
  height: number
}

export const PACK_DRAFT = 0
export const PACK_PUBLISHED = 1
export const PACK_HIDDEN = 2

export const RATING_SFW = 0
export const RATING_NSFW = 1

export const MAX_STICKERS_PER_PACK = 80
export const MAX_TAGS_PER_PACK = 10

/**
 * The database key space (zh-cn / zh-tw / ja-jp / en-us) is wider than the
 * three UI locales, so a lookup walks a per-locale preference order and only
 * then falls back to whatever language the row actually has.
 */
const localeOrder: Record<string, string[]> = {
  'zh-cn': ['zh-cn', 'zh-tw', 'ja-jp', 'en-us'],
  en: ['en-us', 'en', 'ja-jp', 'zh-tw', 'zh-cn'],
  ja: ['ja-jp', 'ja', 'en-us', 'zh-tw', 'zh-cn']
}

export const resolveMultilingual = (
  value: MultilingualText | null | undefined,
  locale: string
): string => {
  if (!value) return ''
  const order = localeOrder[locale] ?? ['zh-cn', 'zh-tw', 'ja-jp', 'en-us']
  for (const key of order) {
    const v = value[key]
    if (v) return v
  }
  return Object.values(value).find(Boolean) ?? ''
}

export const localeField = (locale: string): string => {
  if (locale === 'en') return 'en-us'
  if (locale === 'ja') return 'ja-jp'
  return 'zh-cn'
}

/** The languages the editor offers, in the order they are shown. */
export const EDITABLE_LOCALES = ['zh-cn', 'en-us', 'ja-jp'] as const

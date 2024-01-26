export interface KUNStickers {
  // Sticker pack ID
  sid: number
  // Sticker image ID
  pid: number
  // Sticker image source
  src: string
  // Visual Novel name
  game_en: string
  game_zh: string
  game_ja: string
  // Kawaii loli name
  loli_en: string
  loli_zh: string
  loli_ja: string
  // VNDB VisualNovel vid
  vndb: number
  // Sticker description
  describe: string
}

export interface KUNStickersResponseData {
  sid: number
  pid: number
  game: string
  loli: string
  vndb: number
  describe: string
}

export interface KUNStickerData {
  pid: number
  vndb: number
  describe: string
}

export interface KUNStickers {
  // Sticker pack ID
  sid: number
  // Sticker image ID
  pid: number
  // Sticker image source
  src: string
  // Visual Novel name
  game: KunLanguage
  // Kawaii loli name
  loli: KunLanguage
  // VNDB VisualNovel vid
  vndb: number
  // Sticker description
  describe: string
}

export interface KUNStickersResponseData {
  sid: number
  pid: number
  game: KunLanguage
  loli: KunLanguage
  vndb: number
  describe: string
}

export interface KUNStickerData {
  pid: number
  vndb: number
  describe: string
}

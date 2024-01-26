export interface KUNStickers {
  // 贴纸包 ID
  sid: number
  // 图片 ID
  pid: number
  // 图片源
  src: string
  // 游戏名
  game_en: string
  game_zh: string
  game_ja: string
  // 小只可爱软萌妹子名（咳咳咳）
  loli_en: string
  loli_zh: string
  loli_ja: string
  // VNDB 游戏编号
  vndb: number
  // 图片描述
  describe: string
}

export interface KUNStickerData {
  pid: number
  vndb: number
  describe: string
}

import KUNStickerModel from '~/lib/server/models'
import type { KUNStickersResponseData } from '~/types/stickers'

const findOptions = {
  _id: 0,
  sid: 1,
  pid: 1,
  src: 1,
  vndb: 1,
  describe: 1
}

const findOptionsEN = { ...findOptions, game_en: 1, loli_en: 1 }

const findOptionsZH = { ...findOptions, game_zh: 1, loli_zh: 1 }

export const findStickersData = async (local: App.KunLanguage, sid: number) => {
  const queryOptions = local === 'en' ? findOptionsEN : findOptionsZH

  const data = await KUNStickerModel.find({ sid }, queryOptions)
  return data.map((sticker) => ({
    sid: sticker.sid,
    pid: sticker.pid,
    game: sticker.game_en || sticker.game_ja || sticker.game_zh,
    loli: sticker.loli_en || sticker.loli_ja || sticker.loli_zh,
    vndb: sticker.vndb,
    describe: sticker.describe
  }))
}

export const findOneStickerData = async (local: App.KunLanguage, sid: number, pid: number) => {
  const queryOptions = local === 'en' ? findOptionsEN : findOptionsZH

  const data = await KUNStickerModel.findOne({ sid, pid }, queryOptions)

  if (!data) {
    return
  }

  const stickerData: KUNStickersResponseData = {
    sid: data.sid,
    pid: data.pid,
    game: data.game_en || data.game_ja || data.game_zh,
    loli: data.loli_en || data.loli_ja || data.loli_zh,
    vndb: data.vndb,
    describe: data.describe
  }
  return stickerData
}

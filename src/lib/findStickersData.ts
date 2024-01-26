import KUNStickerModel from '~/lib/server/models'

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

import KUNStickerModel from '~/lib/server/models'
import type { KUNStickersResponseData } from '~/types/stickers'

export const findStickersData = async (sid: number) => {
  const data = await KUNStickerModel.find({ sid }).sort({ pid: 1 }).lean()
  return data.map((sticker) => ({
    sid: sticker.sid,
    pid: sticker.pid,
    game: sticker.game,
    loli: sticker.loli,
    vndb: sticker.vndb,
    describe: sticker.describe
  }))
}

export const findOneStickerData = async (sid: number, pid: number) => {
  const data = await KUNStickerModel.findOne({ sid, pid })

  if (!data) {
    return
  }

  const stickerData: KUNStickersResponseData = {
    sid: data.sid,
    pid: data.pid,
    game: data.game,
    loli: data.loli,
    vndb: data.vndb,
    describe: data.describe
  }
  return stickerData
}

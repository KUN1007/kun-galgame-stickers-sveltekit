import { findOneStickerData } from '~/lib/findStickersData'
export const load = async ({ params, depends }) => {
  const { sid, pid } = params
  depends('kun:details')

  const stickerData = await findOneStickerData(parseInt(sid), parseInt(pid))

  return {
    sid,
    pid,
    game: JSON.stringify(stickerData?.game),
    loli: JSON.stringify(stickerData?.loli)
  }
}

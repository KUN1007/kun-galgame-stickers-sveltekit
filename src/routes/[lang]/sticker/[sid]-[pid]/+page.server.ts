import { findOneStickerData } from '~/lib/findStickersData'
export const load = async ({ params, locals, depends }) => {
  const { sid, pid } = params
  depends('kun:details')

  const stickerData = await findOneStickerData(locals.language, parseInt(sid), parseInt(pid))

  return { sid, pid, stickerData }
}

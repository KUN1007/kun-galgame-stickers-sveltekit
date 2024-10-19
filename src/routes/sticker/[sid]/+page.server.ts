import { findStickersData } from '~/lib/findStickersData'
import type { KUNStickersResponseData } from '~/types/stickers'

export const load = async ({ params, depends }) => {
  const { sid } = params
  depends('kun:sticker')

  const stickersData: KUNStickersResponseData[] = await findStickersData(parseInt(sid))

  return { sid, stickersData }
}

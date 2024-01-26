import { findStickersData } from '~/lib/findStickersData'
import type { KUNStickersResponseData } from '~/types/stickers'

export const load = async ({ params, locals, depends }) => {
  const { sid } = params
  depends('kun:sticker')

  const stickersData: KUNStickersResponseData[] = await findStickersData(
    locals.language,
    parseInt(sid)
  )

  return { sid, stickersData }
}

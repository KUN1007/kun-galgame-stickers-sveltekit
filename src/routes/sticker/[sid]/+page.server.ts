import { getImageArray } from '~/lib/home/stickerSource'

export const load = ({ params }) => {
  const { sid } = params

  if (sid && !isNaN(Number(sid))) {
    const stickersArray = getImageArray(Number(sid))
    return { stickers: { stickersArray } }
  }

  return { stickers: {} }
}

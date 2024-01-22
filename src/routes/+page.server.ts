import { kunStickerPack } from '~/lib/home/stickerSource'

export const load = () => {
  return {
    stickers: kunStickerPack()
  }
}

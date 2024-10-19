import mongoose from 'mongoose'

import type { KUNStickers } from '~/types/stickers'

const KUNStickerSchema = new mongoose.Schema<KUNStickers>(
  {
    sid: { type: Number, default: 0 },
    pid: { type: Number, default: 0 },
    src: { type: String, default: '' },
    game: {
      'en-us': { type: String, default: '' },
      'ja-jp': { type: String, default: '' },
      'zh-cn': { type: String, default: '' },
      'zh-tw': { type: String, default: '' }
    },
    loli: {
      'en-us': { type: String, default: '' },
      'ja-jp': { type: String, default: '' },
      'zh-cn': { type: String, default: '' },
      'zh-tw': { type: String, default: '' }
    },
    vndb: { type: Number, default: 0 },
    describe: { type: String, default: '' }
  },
  { timestamps: { createdAt: 'created', updatedAt: 'updated' } }
)

const KUNStickerModel =
  mongoose.models.stickers || mongoose.model<KUNStickers>('stickers', KUNStickerSchema)

export default KUNStickerModel

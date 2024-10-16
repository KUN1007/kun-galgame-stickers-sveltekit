import mongoose from 'mongoose'

import type { KUNStickers } from '~/types/stickers'

const KUNStickerSchema = new mongoose.Schema<KUNStickers>(
  {
    sid: { type: Number, default: 0 },
    pid: { type: Number, default: 0 },
    src: { type: String, default: '' },

    game_en: { type: String, default: '' },
    game_zh: { type: String, default: '' },
    game_ja: { type: String, default: '' },

    loli_en: { type: String, default: '' },
    loli_zh: { type: String, default: '' },
    loli_ja: { type: String, default: '' },

    vndb: { type: Number, default: 0 },
    describe: { type: String, default: '' }
  },
  { timestamps: { createdAt: 'created', updatedAt: 'updated' } }
)

const KUNStickerModel =
  mongoose.models.stickers || mongoose.model<KUNStickers>('stickers', KUNStickerSchema)

export default KUNStickerModel

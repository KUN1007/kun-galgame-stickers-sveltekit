import mongoose from 'mongoose'

const mongoURI = 'mongodb://127.0.0.1:27017/kungalgame_sticker'
mongoose.connect(mongoURI, { useNewUrlParser: true, useUnifiedTopology: true })

const KUNStickerSchema = new mongoose.Schema(
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
const KUNStickerModel = mongoose.model('stickers', KUNStickerSchema)

async function migrateGameData() {
  const games = await KUNStickerModel.find()

  for (let game of games) {
    console.log(`Processing game with sid: ${game.sid}, pid: ${game.pid}`)
    console.log(
      `Old game values: en=${game.game_en}, ja=${game.game_ja}, zh=${game.game_zh}, tw=${game.game_tw}`
    )
    console.log(
      `Old loli values: en=${game.loli_en}, ja=${game.loli_ja}, zh=${game.loli_zh}, tw=${game.loli_tw}`
    )

    await KUNStickerModel.updateOne(
      { sid: game.sid, pid: game.pid },
      {
        $set: {
          game: {
            'en-us': game.game_en || '',
            'ja-jp': game.game_ja || '',
            'zh-cn': game.game_zh || '',
            'zh-tw': game.game_tw || ''
          },
          loli: {
            'en-us': game.loli_en || '',
            'ja-jp': game.loli_ja || '',
            'zh-cn': game.loli_zh || '',
            'zh-tw': game.loli_tw || ''
          }
        },
        $unset: {
          game_en: '',
          game_ja: '',
          game_zh: '',
          loli_en: '',
          loli_ja: '',
          loli_zh: ''
        }
      }
    )
  }
}

migrateGameData().then(() => console.log('Migration completed.'))

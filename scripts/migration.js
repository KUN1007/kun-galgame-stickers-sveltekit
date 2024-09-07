import path from 'path'
import { promises as fs } from 'fs'
import { fileURLToPath } from 'node:url'
import mongoose from 'mongoose'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

// MongoDB connection string
const mongoURI = 'mongodb://127.0.0.1:27017/kungalgame_sticker'
mongoose.connect(mongoURI, { useNewUrlParser: true, useUnifiedTopology: true })

// Define the schema and model
const KUNStickerSchema = new mongoose.Schema(
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
const KUNStickerModel = mongoose.model('stickers', KUNStickerSchema)

// Function to read and parse a JSON file
async function readAndParseJSON(filePath) {
  const data = await fs.readFile(filePath, 'utf8')
  return JSON.parse(data)
}

// Function to migrate data
async function migrateData() {
  const gameFolderPathEn = path.join(__dirname, 'en/game')
  const lassFolderPathEn = path.join(__dirname, 'en/lass')

  const gameFolderPathZh = path.join(__dirname, 'zh/game')
  const lassFolderPathZh = path.join(__dirname, 'zh/lass')

  for (let i = 1; i <= 6; i++) {
    const gameDataEn = await readAndParseJSON(path.join(gameFolderPathEn, `game${i}.json`))
    const lassDataEn = await readAndParseJSON(path.join(lassFolderPathEn, `lass${i}.json`))
    const gameDataZh = await readAndParseJSON(path.join(gameFolderPathZh, `game${i}.json`))
    const lassDataZh = await readAndParseJSON(path.join(lassFolderPathZh, `lass${i}.json`))

    for (const pid in gameDataEn) {
      // if (Object.hasOwnProperty.call(gameData, pid)) {
      const sticker = new KUNStickerModel({
        sid: i,
        pid: parseInt(pid, 10),
        game_en: gameDataEn[pid],
        loli_en: lassDataEn[pid],
        game_zh: gameDataZh[pid],
        loli_zh: lassDataZh[pid]
      })
      await sticker.save()

      // await KUNStickerModel.updateOne(
      //   { sid: i, pid },
      //   {
      //     game_zh: gameData[pid],
      //     loli_zh: lassData[pid]
      //   }
      // )
      // }
    }
  }
}

// Run the migration
migrateData()
  .then(() => console.log('Migration completed.'))
  .catch((error) => console.error('Migration failed:', error))

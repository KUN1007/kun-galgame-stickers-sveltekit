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
  // const gameFolderPath = path.join(__dirname, 'en/game')
  // const lassFolderPath = path.join(__dirname, 'en/lass')

  const gameFolderPath = path.join(__dirname, 'zh/game')
  const lassFolderPath = path.join(__dirname, 'zh/lass')

  for (let i = 1; i <= 6; i++) {
    const gameData = await readAndParseJSON(path.join(gameFolderPath, `game${i}.json`))
    const lassData = await readAndParseJSON(path.join(lassFolderPath, `lass${i}.json`))

    for (const pid in gameData) {
      if (Object.hasOwnProperty.call(gameData, pid)) {
        // const sticker = new KUNStickerModel({
        //   sid: i,
        //   pid: parseInt(pid, 10),
        //   game_en: gameData[pid],
        //   loli_en: lassData[pid]
        // })
        // await sticker.save()

        await KUNStickerModel.updateOne(
          { sid: i, pid },
          {
            game_zh: gameData[pid],
            loli_zh: lassData[pid]
          }
        )
      }
    }
  }
}

// Run the migration
migrateData()
  .then(() => console.log('Migration completed.'))
  .catch((error) => console.error('Migration failed:', error))

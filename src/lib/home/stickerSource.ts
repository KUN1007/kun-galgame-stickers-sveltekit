import fs from 'fs'
import path from 'path'
import type { KUNStickers } from '~/types/stickers'

const readFileName = (folder: string, folderPath: string, sid: number): KUNStickers[] => {
  const files = fs.readdirSync(folderPath)
  const stickers: KUNStickers[] = files.map((file) => {
    const parts = file.split('-')
    const pid = parseInt(parts[0])
    const loli = parts[1]
    const game = parts[2]
    const describe = parts[3]
    const vndb = parseInt(parts[4])
    return {
      sid: sid,
      pid: pid,
      src: path.join('stickers', folder, file),
      game: game,
      loli: loli,
      vndb: vndb,
      describe: describe
    }
  })

  stickers.sort((a, b) => a.pid - b.pid)
  return stickers
}

export const getImageArray = (sid: number): KUNStickers[] => {
  const stickersDir = path.join(process.cwd(), 'static', 'stickers')
  const folderName = `KUNgal${sid}`
  const folderPath = path.join(stickersDir, folderName)
  if (fs.existsSync(folderPath) && fs.statSync(folderPath).isDirectory()) {
    return readFileName(folderName, folderPath, sid)
  }
  return []
}

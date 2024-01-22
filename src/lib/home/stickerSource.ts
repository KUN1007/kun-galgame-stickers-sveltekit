import fs from 'fs'
import path from 'path'
import type { KUNStickers } from '~/types/stickers'

const readFileName = (folder: string, folderPath: string): KUNStickers[] => {
  const sid = parseInt(folder.split('KUNgal')[1])
  const files = fs.readdirSync(folderPath)
  return files.map((file) => {
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
}

const generateArray = (dir: string): KUNStickers[] => {
  const folders = fs.readdirSync(dir)
  return folders.flatMap((folder) => {
    const folderPath = path.join(dir, folder)
    const stat = fs.statSync(folderPath)
    if (stat.isDirectory()) {
      return readFileName(folder, folderPath)
    }
    return []
  })
}

export const getStickersPath = () => {
  return `${process.cwd()}/static/stickers`
}

export const kunStickerPack = () => {
  const stickersPath = getStickersPath()
  return generateArray(stickersPath)
}

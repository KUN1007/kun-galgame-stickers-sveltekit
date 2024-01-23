import fs from 'node:fs'
import path from 'node:path'
import os from 'node:os'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const stickersDir = path.join(__dirname, '..', 'static', 'stickers')
const jsonDir = path.join(__dirname, '..', 'static', 'json')
const kunGalPrefix = 'KUNgal'

if (!fs.existsSync(jsonDir)) {
  fs.mkdirSync(jsonDir, { recursive: true })
}

const directories = fs.readdirSync(stickersDir).filter((dir) => {
  const dirPath = path.join(stickersDir, dir)
  return fs.statSync(dirPath).isDirectory() && dir.startsWith(kunGalPrefix)
})

directories.forEach((dir) => {
  const dirPath = path.join(stickersDir, dir)
  const files = fs.readdirSync(dirPath)

  const imageData = files.map((file) => {
    let srcPath = path.join('stickers', dir, file)
    if (os.platform() === 'win32') {
      srcPath = srcPath.replace(/\\/g, '/')
    }
    return {
      name: file,
      src: srcPath
    }
  })

  const jsonContent = JSON.stringify(imageData, null, 2)

  const jsonFilePath = path.join(jsonDir, `${dir}.json`)
  fs.writeFileSync(jsonFilePath, jsonContent)
})

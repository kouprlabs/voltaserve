import fs from 'fs'
import path from 'path'
import { rimrafSync } from 'rimraf'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

function cleanup(keepPathsArray: string[]) {
  const distDir = path.join(__dirname, '..', 'dist')
  let rootItems: string[] = []
  try {
    rootItems = fs.readdirSync(distDir)
  } catch (err) {
    console.error(err.message)
    return
  }
  rootItems.forEach((item) => {
    const itemPath = path.join(distDir, item)
    if (!keepPathsArray.includes(item)) {
      try {
        rimrafSync(itemPath, { preserveRoot: false })
      } catch (err) {
        console.error(err.message)
      }
    }
  })
}

cleanup([
  'main.js',
  'main.js.map',
  'main.css',
  'module.js',
  'module.js.map',
  'module.css',
  'types.d.ts',
])

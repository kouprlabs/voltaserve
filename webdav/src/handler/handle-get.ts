import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import os from 'os'
import path from 'path'
import { v4 as uuidv4 } from 'uuid'
import { FileAPI } from '@/client/api'
import { Token } from '@/client/idp'

/*
  This method retrieves the content of a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Create a read stream from the file and pipe it to the response stream.
  - Set the response status code to 200 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
async function handleGet(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  try {
    const api = new FileAPI(token)
    const list = await api.listByPath(req.url)

    const outputPath = path.join(os.tmpdir(), uuidv4())

    api.downloadOriginal(list[0], (response) => {
      const ws = fs.createWriteStream(outputPath)
      response.pipe(ws)
      ws.on('finish', () => {
        ws.close()
        const rs = fs.createReadStream(outputPath)
        rs.on('error', (error) => {
          console.error(error)
          res.statusCode = 500
          res.end()
        })
        rs.on('end', () => {
          fs.rmSync(outputPath)
        })
        rs.pipe(res)
      })
    })
  } catch {
    res.statusCode = 500
    res.end()
  }
}

export default handleGet

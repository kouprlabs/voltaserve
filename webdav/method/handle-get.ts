import fs from 'fs'
import { IncomingMessage, ServerResponse, get } from 'http'
import os from 'os'
import path from 'path'
import { v4 as uuidv4 } from 'uuid'
import { File } from '@/api/file'
import { Token } from '@/api/token'
import { API_URL } from '@/config/config'

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
    const result = await fetch(`${API_URL}/v1/files/list?path=${req.url}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token.access_token}`,
      },
    })
    const file: File = (await result.json())[0]
    const filePath = path.join(os.tmpdir(), uuidv4())
    get(
      `${API_URL}/v1/files/${file.id}/original${file.original.extension}?access_token=${token.access_token}`,
      (response) => {
        const ws = fs.createWriteStream(filePath)
        response.pipe(ws)
        ws.on('finish', () => {
          ws.close()
          const rs = fs.createReadStream(filePath)
          rs.on('error', (error) => {
            console.error(error)
            res.statusCode = 500
            res.end()
          })
          rs.on('end', () => {
            fs.rmSync(filePath)
          })
          rs.pipe(res)
        })
      }
    )
  } catch {
    res.statusCode = 500
    res.end()
  }
}

export default handleGet

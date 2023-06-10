import fs, { createReadStream, rmSync } from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import os from 'os'
import path from 'path'
import { v4 as uuidv4 } from 'uuid'
import { FileAPI } from '@/client/api'
import { Token } from '@/client/idp'
import { handleException } from '@/infra/error'

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
    const file = await api.getByPath(decodeURI(req.url))

    /* TODO: This should be optimized for the case when there is a range header,
       only a partial file should be fetched, here we are fetching the whole file
       which is not ideal. */
    const outputPath = path.join(os.tmpdir(), uuidv4())
    await api.downloadOriginal(file, outputPath)

    const stat = fs.statSync(outputPath)
    const rangeHeader = req.headers.range
    if (rangeHeader) {
      const [start, end] = rangeHeader.replace(/bytes=/, '').split('-')
      const rangeStart = parseInt(start, 10) || 0
      const rangeEnd = parseInt(end, 10) || stat.size - 1
      const chunkSize = rangeEnd - rangeStart + 1
      res.writeHead(206, {
        'Content-Range': `bytes ${rangeStart}-${rangeEnd}/${stat.size}`,
        'Accept-Ranges': 'bytes',
        'Content-Length': chunkSize.toString(),
        'Content-Type': 'application/octet-stream',
      })
      createReadStream(outputPath, {
        start: rangeStart,
        end: rangeEnd,
      })
        .pipe(res)
        .on('finish', () => rmSync(outputPath))
    } else {
      res.writeHead(200, {
        'Content-Length': stat.size.toString(),
        'Content-Type': 'application/octet-stream',
      })
      createReadStream(outputPath)
        .pipe(res)
        .on('finish', () => rmSync(outputPath))
    }
  } catch (err) {
    handleException(err, res)
  }
}

export default handleGet

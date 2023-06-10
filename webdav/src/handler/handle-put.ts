import fs from 'fs'
import { readFile } from 'fs/promises'
import { IncomingMessage, ServerResponse } from 'http'
import os from 'os'
import path from 'path'
import { v4 as uuidv4 } from 'uuid'
import { Token } from '@/client/idp'
import { FileAPI } from '@/client/api'
import { handleException } from '@/infra/error'

/*
  This method creates or updates a resource with the provided content.

  Example implementation:

  - Extract the file path from the URL.
  - Create a write stream to the file.
  - Listen for the data event to write the incoming data to the file.
  - Listen for the end event to indicate the completion of the write stream.
  - Set the response status code to 201 if created or 204 if updated.
  - Return the response.
 */
async function handlePut(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  const api = new FileAPI(token)
  try {
    const directory = await api.getByPath(decodeURI(path.dirname(req.url)))
    const outputPath = path.join(os.tmpdir(), uuidv4())
    const ws = fs.createWriteStream(outputPath)
    req.pipe(ws)
    ws.on('error', (err) => {
      console.error(err)
      res.statusCode = 500
      res.end()
    })
    ws.on('finish', async () => {
      try {
        res.statusCode = 201
        res.end()

        const blob = new Blob([await readFile(outputPath)])

        await api.upload({
          workspaceId: directory.workspaceId,
          parentId: directory.id,
          name: decodeURIComponent(path.basename(req.url)),
          blob,
        })

        /* Delete existing file (simulate an overwrite) */
        const existingFile = await api.getByPath(decodeURI(req.url))
        await api.delete(existingFile.id)
      } catch (err) {
        handleException(err, res)
      } finally {
        fs.rmSync(outputPath)
      }
    })
  } catch (err) {
    handleException(err, res)
  }
}

export default handlePut

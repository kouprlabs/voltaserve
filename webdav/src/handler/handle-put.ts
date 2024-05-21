import fs from 'fs'
import { readFile } from 'fs/promises'
import { IncomingMessage, ServerResponse } from 'http'
import fsPromises from 'node:fs/promises'
import os from 'os'
import path from 'path'
import { v4 as uuidv4 } from 'uuid'
import { FileAPI, FileType } from '@/client/api'
import { Token } from '@/client/idp'
import {
  isMicrosoftOfficeLockFile,
  isOpenOfficeOfficeLockFile,
} from '@/helper/office-lock-files'
import { handleError } from '@/infra/error'

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
  token: Token,
) {
  const name = decodeURIComponent(path.basename(req.url))
  if (isMicrosoftOfficeLockFile(name) || isOpenOfficeOfficeLockFile(name)) {
    res.statusCode = 200
    res.end()
  } else {
    const api = new FileAPI(token)
    try {
      const directory = await api.getByPath(
        decodeURIComponent(path.dirname(req.url)),
      )
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
          const blob = new Blob([await readFile(outputPath)])
          const name = decodeURIComponent(path.basename(req.url))
          try {
            const existingFile = await api.getByPath(
              decodeURIComponent(req.url),
            )
            await api.patch({
              id: existingFile.id,
              blob,
              name,
            })
            res.statusCode = 201
            res.end()
          } catch {
            await api.create({
              type: FileType.File,
              workspaceId: directory.workspaceId,
              parentId: directory.id,
              blob,
              name,
            })
            res.statusCode = 201
            res.end()
          }
        } catch (err) {
          handleError(err, res)
        } finally {
          await fsPromises.rm(outputPath)
        }
      })
    } catch (err) {
      handleError(err, res)
    }
  }
}

export default handlePut

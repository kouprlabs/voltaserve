import fs from 'fs'
import { readFile } from 'fs/promises'
import { IncomingMessage, ServerResponse } from 'http'
import os from 'os'
import path from 'path'
import { v4 as uuidv4 } from 'uuid'
import { File } from '@/api/file'
import { Token } from '@/api/token'
import { API_URL } from '@/config/config'

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
  let directory: File
  try {
    const result = await fetch(
      `${API_URL}/v1/files/get?path=${path.dirname(req.url)}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token.access_token}`,
        },
      }
    )
    directory = (await result.json()) as File

    const filePath = path.join(os.tmpdir(), uuidv4())
    const ws = fs.createWriteStream(filePath)
    ws.on('error', (error) => {
      console.error(error)
      res.statusCode = 500
      res.end()
    })
    req.on('data', (chunk) => {
      ws.write(chunk)
    })
    req.on('end', async () => {
      ws.end()

      const params = new URLSearchParams({
        workspace_id: directory.workspaceId,
      })
      params.append('parent_id', directory.id)

      const formData = new FormData()
      const blob = new Blob([await readFile(filePath)])
      formData.set('file', blob, path.basename(req.url))

      try {
        await fetch(`${API_URL}/v1/files?${params}`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token.access_token}`,
          },
          body: formData,
        })
        res.statusCode = 201
      } catch (err) {
        console.error(err)
        res.statusCode = 500
      }
      fs.rmSync(filePath)
      res.end()
    })
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
    return
  }
}

export default handlePut

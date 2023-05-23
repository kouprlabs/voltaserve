import fs from 'fs'
import { readFile } from 'fs/promises'
import { IncomingMessage, ServerResponse } from 'http'
import os from 'os'
import path from 'path'
import { v4 as uuidv4 } from 'uuid'
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
  try {
    /* Delete existing file (simulate an overwrite) */
    const result = await fetch(`${API_URL}/v1/files/get?path=${req.url}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    const file = await result.json()
    await fetch(`${API_URL}/v1/files/${file.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token.access_token}`,
        'Content-Type': 'application/json',
      },
    })
  } catch (err) {}
  try {
    const result = await fetch(
      `${API_URL}/v1/files/get?path=${path.dirname(req.url)}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token.access_token}`,
          'Content-Type': 'application/json',
        },
      }
    )
    const directory = await result.json()

    const filePath = path.join(os.tmpdir(), uuidv4())
    const ws = fs.createWriteStream(filePath)
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

        const params = new URLSearchParams({
          workspace_id: directory.workspaceId,
        })
        params.append('parent_id', directory.id)

        const formData = new FormData()
        const blob = new Blob([await readFile(filePath)])
        formData.set('file', blob, decodeURIComponent(path.basename(req.url)))

        await fetch(`${API_URL}/v1/files?${params}`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token.access_token}`,
          },
          body: formData,
        })
      } catch (err) {
        console.error(err)
        res.statusCode = 500
        res.end()
      } finally {
        fs.rmSync(filePath)
      }
    })
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
    return
  }
}

export default handlePut

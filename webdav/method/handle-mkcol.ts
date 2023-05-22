import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { File } from '@/api/file'
import { Token } from '@/api/token'
import { API_URL } from '@/config/config'
import { getFilePath } from '@/infra/path'

/*
  This method creates a new collection (directory) at the specified URL.

  Example implementation:

  - Extract the directory path from the URL.
  - Use fs.mkdir() to create the directory.
  - Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
  - Return the response.
 */
async function handleMkcol(
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
          'Content-Type': 'application/json',
        },
      }
    )
    directory = await result.json()
    await fetch(`${API_URL}/v1/files/create_folder`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        workspaceId: directory.workspaceId,
        parentId: directory.id,
        name: decodeURIComponent(path.basename(req.url)),
      }),
    })
    res.statusCode = 201
    res.end()
    console.log(`handleMkcol: ${req.url} [OK]`)
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handleMkcol

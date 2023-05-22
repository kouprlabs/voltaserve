import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import { File, FileType } from '@/api/file'
import { Token } from '@/api/token'
import { API_URL } from '@/config/config'
import { getFilePath } from '@/infra/path'

/*
  This method is similar to GET but only retrieves the metadata of a resource, without returning the actual content.

  Example implementation:

  - Extract the file path from the URL.
  - Retrieve the file metadata using fs.stat().
  - Set the response status code to 200 if successful or an appropriate error code if the file is not found.
  - Set the Content-Length header with the file size.
  - Return the response.
*/
async function handleHead(
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
    if (file.type === FileType.File) {
      res.statusCode = 200
      res.setHeader('Content-Length', file.original.size)
      res.end()
    } else {
      res.statusCode = 200
      res.end()
    }
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handleHead

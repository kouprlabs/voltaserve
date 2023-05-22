import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import { Token } from '@/api/token'
import { getFilePath } from '@/infra/path'

/*
  This method deletes a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Use fs.unlink() to delete the file.
  - Set the response status code to 204 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
function handleDelete(req: IncomingMessage, res: ServerResponse, token: Token) {
  const filePath = getFilePath(req.url)
  fs.rm(filePath, { recursive: true }, (error) => {
    if (error) {
      console.error(error)
      if (error.code === 'ENOENT') {
        res.statusCode = 404
      } else {
        res.statusCode = 500
      }
    } else {
      res.statusCode = 204
    }
    res.end()
  })
}

export default handleDelete

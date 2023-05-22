import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import { Token } from '@/api/token'
import { getFilePath } from '@/infra/path'

/*
  This method creates a new collection (directory) at the specified URL.

  Example implementation:

  - Extract the directory path from the URL.
  - Use fs.mkdir() to create the directory.
  - Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
  - Return the response.
 */
function handleMkcol(req: IncomingMessage, res: ServerResponse, token: Token) {
  const filePath = getFilePath(req.url)
  fs.mkdir(filePath, (error) => {
    if (error) {
      console.error(error)
      res.statusCode = 500
    } else {
      res.statusCode = 201
    }
    res.end()
  })
}

export default handleMkcol

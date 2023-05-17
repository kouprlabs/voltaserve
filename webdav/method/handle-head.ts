import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
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
function handleHead(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  fs.stat(filePath, (error, stats) => {
    if (error) {
      console.error(error)
      if (error.code === 'ENOENT') {
        res.statusCode = 404
      } else {
        res.statusCode = 500
      }
      res.end()
    } else {
      res.statusCode = 200
      res.setHeader('Content-Length', stats.size)
      res.end()
    }
  })
}

export default handleHead

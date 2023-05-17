import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import { getDestinationPath, getFilePath } from '@/infra/path'

/*
  This method moves or renames a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.rename() to move or rename the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
function handleMove(req: IncomingMessage, res: ServerResponse) {
  const sourcePath = getFilePath(req.url)
  const destinationPath = getDestinationPath(req)
  fs.rename(sourcePath, destinationPath, (error) => {
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

export default handleMove

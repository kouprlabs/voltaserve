import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import { getFilePath } from '@/helper/path'

/*
  This method retrieves the content of a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Create a read stream from the file and pipe it to the response stream.
  - Set the response status code to 200 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
function handleGet(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  const fileStream = fs.createReadStream(filePath)
  fileStream.on('error', (error) => {
    console.error(error)
    res.statusCode = 500
    res.end()
  })
  fileStream.pipe(res)
}

export default handleGet

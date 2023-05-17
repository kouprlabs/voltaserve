import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import { getFilePath } from '@/infra/path'

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
function handlePut(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  const fileStream = fs.createWriteStream(filePath)
  fileStream.on('error', (error) => {
    console.error(error)
    res.statusCode = 500
    res.end()
  })
  req.on('data', (chunk) => {
    fileStream.write(chunk)
  })
  req.on('end', () => {
    fileStream.end()
    res.statusCode = 201
    res.end()
  })
}

export default handlePut

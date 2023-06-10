import { IncomingMessage, ServerResponse } from 'http'
import { FileAPI, FileType } from '@/client/api'
import { Token } from '@/client/idp'
import { handleException } from '@/infra/error'

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
    const file = await new FileAPI(token).getByPath(decodeURIComponent(req.url))
    if (file.type === FileType.File) {
      res.statusCode = 200
      res.setHeader('Content-Length', file.original.size)
      res.end()
    } else {
      res.statusCode = 200
      res.end()
    }
  } catch (err) {
    handleException(err, res)
  }
}

export default handleHead

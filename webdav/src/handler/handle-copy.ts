import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { getTargetPath } from '@/helper/path'
import { FileAPI } from '@/client/api'
import { Token } from '@/client/idp'
import { handleException } from '@/infra/error'

/*
  This method copies a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.copyFile() to copy the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
async function handleCopy(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  try {
    const api = new FileAPI(token)
    const sourceFile = await api.getByPath(decodeURI(req.url))
    const targetFile = await api.getByPath(
      decodeURI(path.dirname(getTargetPath(req)))
    )

    if (sourceFile.workspaceId !== targetFile.workspaceId) {
      res.statusCode = 400
      res.end()
      return
    }

    const clones = await api.copy(targetFile.id, { ids: [sourceFile.id] })
    await api.rename(clones[0].id, {
      name: decodeURIComponent(path.basename(getTargetPath(req))),
    })

    res.statusCode = 204
    res.end()
  } catch (err) {
    handleException(err, res)
  }
}

export default handleCopy

import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { FileAPI } from '@/client/api'
import { Token } from '@/client/idp'
import { getTargetPath } from '@/infra/path'

/*
  This method moves or renames a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.rename() to move or rename the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
async function handleMove(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  try {
    const sourcePath = decodeURI(req.url)
    const targetPath = getTargetPath(req)

    const api = new FileAPI(token)
    const sourceFile = await api.getByPath(req.url)
    const targetFile = await api.getByPath(path.dirname(getTargetPath(req)))

    if (sourceFile.workspaceId !== targetFile.workspaceId) {
      res.statusCode = 400
      res.end()
      return
    }

    if (
      sourcePath.split('/').length === targetPath.split('/').length &&
      path.dirname(sourcePath) === path.dirname(targetPath)
    ) {
      await api.rename(sourceFile.id, { name: path.basename(targetPath) })
    } else {
      await api.move(targetFile.id, { ids: [sourceFile.id] })
    }

    res.statusCode = 204
    res.end()
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handleMove

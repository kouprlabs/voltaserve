import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { getTargetPath } from '@/helper/path'
import { FileAPI, geEditorPermission } from '@/client/api'
import { Token } from '@/client/idp'

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

    if (
      !geEditorPermission(sourceFile.permission) ||
      !geEditorPermission(targetFile.permission)
    ) {
      res.statusCode = 401
      res.end()
      return
    }

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
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handleCopy

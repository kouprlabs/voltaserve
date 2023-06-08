import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { File, FileAPI } from '@/client/api'
import { Token } from '@/client/idp'

/*
  This method creates a new collection (directory) at the specified URL.

  Example implementation:

  - Extract the directory path from the URL.
  - Use fs.mkdir() to create the directory.
  - Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
  - Return the response.
 */
async function handleMkcol(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  let directory: File
  try {
    const api = new FileAPI(token)
    directory = await api.getByPath(decodeURI(path.dirname(req.url)))
    await api.createFolder({
      workspaceId: directory.workspaceId,
      parentId: directory.id,
      name: decodeURIComponent(path.basename(req.url)),
    })

    res.statusCode = 201
    res.end()
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handleMkcol

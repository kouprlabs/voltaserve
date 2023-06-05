import { IncomingMessage, ServerResponse } from 'http'
import { FileAPI } from '@/client/api'
import { Token } from '@/client/idp'

/*
  This method deletes a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Use fs.unlink() to delete the file.
  - Set the response status code to 204 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
async function handleDelete(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  try {
    const api = new FileAPI(token)
    const file = await api.getByPath(req.url)
    await api.delete(file.id)

    res.statusCode = 204
    res.end()
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handleDelete

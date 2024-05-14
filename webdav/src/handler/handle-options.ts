import { IncomingMessage, ServerResponse } from 'http'

/*
  This method should respond with the allowed methods and capabilities of the server.

  Example implementation:

  - Set the response status code to 200.
  - Set the Allow header to specify the supported methods, such as OPTIONS, GET, PUT, DELETE, etc.
  - Return the response.
 */
async function handleOptions(_: IncomingMessage, res: ServerResponse) {
  res.statusCode = 200
  res.setHeader(
    'Allow',
    'OPTIONS, GET, HEAD, PUT, DELETE, MKCOL, COPY, MOVE, PROPFIND, PROPPATCH',
  )
  res.end()
}

export default handleOptions

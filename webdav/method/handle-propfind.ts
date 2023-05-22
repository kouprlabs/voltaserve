import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { File, FileType } from '@/api/file'
import { Token } from '@/api/token'
import { API_URL } from '@/config/config'

/*
  This method retrieves properties and metadata of a resource.

  Example implementation:

  - Extract the file path from the URL.
  - Use fs.stat() to retrieve the file metadata.
  - Format the response body in the desired XML format with the properties and metadata.
  - Set the response status code to 207 if successful or an appropriate error code if the file is not found or encountered an error.
  - Set the Content-Type header to indicate the XML format.
  - Return the response.
 */
async function handlePropfind(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  try {
    const result = await fetch(`${API_URL}/v1/files/get?path=${req.url}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    const file = await result.json()
    if (file.type === FileType.File) {
      const responseXml = `
      <D:multistatus xmlns:D="DAV:">
        <D:response>
          <D:href>${encodeURIComponent(file.name)}</D:href>
          <D:propstat>
            <D:prop>
              <D:resourcetype></D:resourcetype>
            </D:prop>
            <D:status>HTTP/1.1 200 OK</D:status>
          </D:propstat>
        </D:response>
      </D:multistatus>`
      res.statusCode = 207
      res.setHeader('Content-Type', 'application/xml; charset=utf-8')
      res.end(responseXml)
    } else if (file.type === FileType.Folder) {
      const result = await fetch(`${API_URL}/v1/files/list?path=${req.url}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token.access_token}`,
          'Content-Type': 'application/json',
        },
      })
      const files: File[] = await result.json()
      const responseXml = `
        <D:multistatus xmlns:D="DAV:">
          <D:response>
            <D:href>${req.url}</D:href>
            <D:propstat>
              <D:prop>
                <D:resourcetype><D:collection/></D:resourcetype>
              </D:prop>
              <D:status>HTTP/1.1 200 OK</D:status>
            </D:propstat>
          </D:response>
          ${files
            .map((item) => {
              return `
                <D:response>
                  <D:href>${path.join(
                    req.url,
                    encodeURIComponent(item.name)
                  )}</D:href>
                  <D:propstat>
                    <D:prop>
                      <D:resourcetype>${
                        item.type === FileType.Folder ? '<D:collection/>' : ''
                      }</D:resourcetype>
                    </D:prop>
                    <D:status>HTTP/1.1 200 OK</D:status>
                  </D:propstat>
                </D:response>
              `
            })
            .join('')}
        </D:multistatus>`
      res.statusCode = 207
      res.setHeader('Content-Type', 'application/xml; charset=utf-8')
      res.end(responseXml)
    }
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handlePropfind

import fs from 'fs'
import { IncomingMessage, ServerResponse, get } from 'http'
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
  const result = await fetch(`${API_URL}/v1/files/list?path=${req.url}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token.access_token}`,
    },
  })
  const files: File[] = await result.json()
  if (files.length === 1 && files[0].type === FileType.File) {
    const responseXml = `
    <D:multistatus xmlns:D="DAV:">
      <D:response>
        <D:href>${encodeURIComponent(files[0].name)}</D:href>
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
  } else {
    const responseXml = `
    <D:multistatus xmlns:D="DAV:">
      ${files
        .map((item) => {
          return `
            <D:response>
              <D:href>${encodeURIComponent(item.name)}</D:href>
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
}

export default handlePropfind

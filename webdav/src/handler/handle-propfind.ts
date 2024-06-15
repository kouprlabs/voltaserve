import { IncomingMessage, ServerResponse } from 'http'
import { FileAPI, FileType } from '@/client/api'
import { Token } from '@/client/idp'
import { handleError } from '@/infra/error'

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
  token: Token,
) {
  try {
    const api = new FileAPI(token)
    const file = await api.getByPath(decodeURIComponent(req.url))
    if (file.type === FileType.File) {
      const responseXml = `
        <D:multistatus xmlns:D="DAV:">
          <D:response>
            <D:href>${encodeURIComponent(file.name)}</D:href>
            <D:propstat>
              <D:prop>
                <D:resourcetype></D:resourcetype>
                ${
                  file.snapshot.original
                    ? `<D:getcontentlength>${file.snapshot.original.size}</D:getcontentlength>`
                    : ''
                }
                <D:creationdate>${new Date(
                  file.createTime,
                ).toUTCString()}</D:creationdate>
                <D:getlastmodified>${new Date(
                  file.updateTime,
                ).toUTCString()}</D:getlastmodified>
              </D:prop>
              <D:status>HTTP/1.1 200 OK</D:status>
            </D:propstat>
          </D:response>
        </D:multistatus>`
      res.statusCode = 207
      res.setHeader('Content-Type', 'application/xml; charset=utf-8')
      res.end(responseXml)
    } else if (file.type === FileType.Folder) {
      const list = await api.listByPath(decodeURIComponent(req.url))
      const responseXml = `
        <D:multistatus xmlns:D="DAV:">
          <D:response>
            <D:href>${req.url}</D:href>
            <D:propstat>
              <D:prop>
                <D:resourcetype><D:collection/></D:resourcetype>
                <D:getcontentlength>0</D:getcontentlength>
                <D:getlastmodified>${new Date(
                  file.updateTime,
                ).toUTCString()}</D:getlastmodified>
                <D:creationdate>${new Date(
                  file.createTime,
                ).toUTCString()}</D:creationdate>
              </D:prop>
              <D:status>HTTP/1.1 200 OK</D:status>
            </D:propstat>
          </D:response>
          ${list
            .map((item) => {
              return `
                <D:response>
                  <D:href>${req.url}${encodeURIComponent(item.name)}</D:href>
                  <D:propstat>
                    <D:prop>
                      <D:resourcetype>${
                        item.type === FileType.Folder ? '<D:collection/>' : ''
                      }</D:resourcetype>
                      ${
                        item.type === FileType.File && item.snapshot.original
                          ? `<D:getcontentlength>${item.snapshot.original.size}</D:getcontentlength>`
                          : ''
                      }
                      <D:getlastmodified>${new Date(
                        item.updateTime,
                      ).toUTCString()}</D:getlastmodified>
                      <D:creationdate>${new Date(
                        item.createTime,
                      ).toUTCString()}</D:creationdate>
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
    handleError(err, res)
  }
}

export default handlePropfind

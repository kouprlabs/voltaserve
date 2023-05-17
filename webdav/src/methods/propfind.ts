import fs from 'fs'
import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { getFilePath } from '@/helper/path'

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
function handlePropfind(req: IncomingMessage, res: ServerResponse) {
  const url = req.url
  const filePath = getFilePath(url)
  fs.stat(filePath, (error, stats) => {
    if (error) {
      console.error(error)
      if (error.code === 'ENOENT') {
        res.statusCode = 404
      } else {
        res.statusCode = 500
      }
      res.end()
      return
    }
    const isDirectory = stats.isDirectory()
    if (isDirectory) {
      fs.readdir(filePath, (error, files) => {
        if (error) {
          res.statusCode = 500
          res.end()
          return
        }
        const responseXml = buildDirectoryPropfindResponse(filePath, url, files)
        res.statusCode = 207
        res.setHeader('Content-Type', 'application/xml; charset=utf-8')
        res.end(responseXml)
      })
    } else {
      const responseXml = buildFilePropfindResponse(url)
      res.statusCode = 207
      res.setHeader('Content-Type', 'application/xml; charset=utf-8')
      res.end(responseXml)
    }
  })
}

function buildDirectoryPropfindResponse(
  directoryPath: string,
  url: string,
  items: string[]
) {
  return `
    <D:multistatus xmlns:D="DAV:">
      <D:response>
        <D:href>${encodeURIComponent(url)}</D:href>
        <D:propstat>
          <D:prop>
            <D:resourcetype>
              <D:collection/>
            </D:resourcetype>
          </D:prop>
          <D:status>HTTP/1.1 200 OK</D:status>
        </D:propstat>
      </D:response>
      ${items
        .map((item) => {
          const stat = fs.statSync(path.join(directoryPath, item))
          return `
            <D:response>
              <D:href>${encodeURIComponent(item)}</D:href>
              <D:propstat>
                <D:prop>
                  <D:resourcetype>${
                    stat.isDirectory() ? '<D:collection/>' : ''
                  }</D:resourcetype>
                </D:prop>
                <D:status>HTTP/1.1 200 OK</D:status>
              </D:propstat>
            </D:response>
          `
        })
        .join('')}
    </D:multistatus>`
}

function buildFilePropfindResponse(filePath: string) {
  return `
    <D:multistatus xmlns:D="DAV:">
      <D:response>
        <D:href>${encodeURIComponent(filePath)}</D:href>
        <D:propstat>
          <D:prop>
            <D:resourcetype></D:resourcetype>
          </D:prop>
          <D:status>HTTP/1.1 200 OK</D:status>
        </D:propstat>
      </D:response>
    </D:multistatus>`
}

export default handlePropfind

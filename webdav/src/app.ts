import fs from 'fs'
import { createServer, IncomingMessage, ServerResponse } from 'http'
import passport from 'passport'
import { BasicStrategy } from 'passport-http'
import path from 'path'

const DATA_DIRECTORY =
  path.join(__dirname, '..', 'data') || process.env.DATA_DIRECTORY

type User = {
  id: number
  username: string
  password: string
}

const users: User[] = [{ id: 1, username: 'admin', password: 'admin' }]

passport.use(
  new BasicStrategy((username, password, done) => {
    const user = users.find((u) => u.username === username)
    if (!user) {
      return done(null, false)
    }
    if (password !== user.password) {
      return done(null, false)
    }
    return done(null, user)
  })
)

const server = createServer((req: IncomingMessage, res: ServerResponse) => {
  passport.authenticate(
    'basic',
    { session: false },
    (err: Error, user: User) => {
      if (err || !user) {
        res.statusCode = 401
        res.setHeader('WWW-Authenticate', 'Basic realm="WebDAV Server"')
        res.end()
        return
      }
      const method = req.method
      switch (method) {
        case 'OPTIONS':
          handleOptions(req, res)
          break
        case 'GET':
          handleGet(req, res)
          break
        case 'HEAD':
          handleHead(req, res)
          break
        case 'PUT':
          handlePut(req, res)
          break
        case 'DELETE':
          handleDelete(req, res)
          break
        case 'MKCOL':
          handleMkcol(req, res)
          break
        case 'COPY':
          handleCopy(req, res)
          break
        case 'MOVE':
          handleMove(req, res)
          break
        case 'PROPFIND':
          handlePropfind(req, res)
          break
        case 'PROPPATCH':
          handleProppatch(req, res)
          break
        default:
          res.statusCode = 501
          res.end()
      }
    }
  )(req, res)
})

/*
  This method should respond with the allowed methods and capabilities of the server.

  Example implementation:

  - Set the response status code to 200.
  - Set the Allow header to specify the supported methods, such as OPTIONS, GET, PUT, DELETE, etc.
  - Return the response.
 */
function handleOptions(_: IncomingMessage, res: ServerResponse) {
  res.statusCode = 200
  res.setHeader(
    'Allow',
    'OPTIONS, GET, HEAD, PUT, DELETE, MKCOL, COPY, MOVE, PROPFIND, PROPPATCH'
  )
  res.end()
}

/*
  This method retrieves the content of a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Create a read stream from the file and pipe it to the response stream.
  - Set the response status code to 200 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
function handleGet(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  const fileStream = fs.createReadStream(filePath)
  fileStream.on('error', (error) => {
    console.error(error)
    res.statusCode = 500
    res.end()
  })
  fileStream.pipe(res)
}

/*
  This method is similar to GET but only retrieves the metadata of a resource, without returning the actual content.

  Example implementation:

  - Extract the file path from the URL.
  - Retrieve the file metadata using fs.stat().
  - Set the response status code to 200 if successful or an appropriate error code if the file is not found.
  - Set the Content-Length header with the file size.
  - Return the response.
*/
function handleHead(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  fs.stat(filePath, (error, stats) => {
    if (error) {
      console.error(error)
      if (error.code === 'ENOENT') {
        res.statusCode = 404
      } else {
        res.statusCode = 500
      }
      res.end()
    } else {
      res.statusCode = 200
      res.setHeader('Content-Length', stats.size)
      res.end()
    }
  })
}

/*
  This method creates or updates a resource with the provided content.

  Example implementation:

  - Extract the file path from the URL.
  - Create a write stream to the file.
  - Listen for the data event to write the incoming data to the file.
  - Listen for the end event to indicate the completion of the write stream.
  - Set the response status code to 201 if created or 204 if updated.
  - Return the response.
 */
function handlePut(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  const fileStream = fs.createWriteStream(filePath)
  fileStream.on('error', (error) => {
    console.error(error)
    res.statusCode = 500
    res.end()
  })
  req.on('data', (chunk) => {
    fileStream.write(chunk)
  })
  req.on('end', () => {
    fileStream.end()
    res.statusCode = 201
    res.end()
  })
}

/*
  This method deletes a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Use fs.unlink() to delete the file.
  - Set the response status code to 204 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
function handleDelete(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  fs.rm(filePath, { recursive: true }, (error) => {
    if (error) {
      console.error(error)
      if (error.code === 'ENOENT') {
        res.statusCode = 404
      } else {
        res.statusCode = 500
      }
    } else {
      res.statusCode = 204
    }
    res.end()
  })
}

/*
  This method creates a new collection (directory) at the specified URL.

  Example implementation:

  - Extract the directory path from the URL.
  - Use fs.mkdir() to create the directory.
  - Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
  - Return the response.
 */
function handleMkcol(req: IncomingMessage, res: ServerResponse) {
  const filePath = getFilePath(req.url)
  fs.mkdir(filePath, (error) => {
    if (error) {
      console.error(error)
      res.statusCode = 500
    } else {
      res.statusCode = 201
    }
    res.end()
  })
}

/*
  This method copies a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.copyFile() to copy the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
function handleCopy(req: IncomingMessage, res: ServerResponse) {
  const sourcePath = getFilePath(req.url)
  const destinationPath = getDestinationPath(req)
  fs.copyFile(sourcePath, destinationPath, (error) => {
    if (error) {
      console.error(error)
      if (error.code === 'ENOENT') {
        res.statusCode = 404
      } else {
        res.statusCode = 500
      }
    } else {
      res.statusCode = 204
    }
    res.end()
  })
}

/*
  This method moves or renames a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.rename() to move or rename the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
function handleMove(req: IncomingMessage, res: ServerResponse) {
  const sourcePath = getFilePath(req.url)
  const destinationPath = getDestinationPath(req)
  fs.rename(sourcePath, destinationPath, (error) => {
    if (error) {
      console.error(error)
      if (error.code === 'ENOENT') {
        res.statusCode = 404
      } else {
        res.statusCode = 500
      }
    } else {
      res.statusCode = 204
    }
    res.end()
  })
}

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

/*
  This method updates the properties of a resource.

  Example implementation:

  - Parse the request body to extract the properties to be updated.
  - Read the existing data from the file.
  - Parse the existing properties.
  - Merge the updated properties with the existing ones.
  - Format the updated properties and store them back in the file.
  - Set the response status code to 204 if successful or an appropriate error code if the file is not found or encountered an error.
  - Return the response.

  In this example implementation, the handleProppatch() method first parses the XML 
  payload containing the properties to be updated. Then, it reads the existing data from the file, 
  parses the existing properties (assuming an XML format), 
  merges the updated properties with the existing ones, and formats 
  the properties back into the desired format (e.g., XML).

  Finally, the updated properties are written back to the file. 
  You can customize the parseProperties() and formatProperties() 
  functions to match the specific property format you are using in your WebDAV server.

  Note that this implementation assumes a simplified example and may require further 
  customization based on your specific property format and requirements.
 */
function handleProppatch(_: IncomingMessage, res: ServerResponse) {
  res.statusCode = 501
  res.end()
}

function getFilePath(url: string) {
  return path.join(DATA_DIRECTORY, decodeURIComponent(url))
}

function getDestinationPath(req: IncomingMessage) {
  const destinationHeader = req.headers.destination as string
  if (!destinationHeader) {
    return null
  }
  // Check if the destination header is a full URL
  if (
    destinationHeader.startsWith('http://') ||
    destinationHeader.startsWith('https://')
  ) {
    const url = new URL(destinationHeader)
    return path.join(DATA_DIRECTORY, decodeURIComponent(url.pathname))
  } else {
    /* Extract the path from the destination header */
    const startIndex =
      destinationHeader.indexOf(req.headers.host) + req.headers.host.length
    const value = destinationHeader.substring(startIndex)
    return path.join(DATA_DIRECTORY, decodeURIComponent(value))
  }
}

const port = 9988 || process.env.PORT
server.listen(port, () => {
  console.log(`WebDAV server is listening on port ${port}`)
})

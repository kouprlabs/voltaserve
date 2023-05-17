import { createServer, IncomingMessage, ServerResponse } from 'http'
import handleCopy from 'methods/copy'
import handleDelete from 'methods/delete'
import handleGet from 'methods/get'
import handleHead from 'methods/head'
import handleMkcol from 'methods/mkcol'
import handleMove from 'methods/move'
import handleOptions from 'methods/options'
import handlePropfind from 'methods/propfind'
import handleProppatch from 'methods/proppatch'
import handlePut from 'methods/put'
import passport from 'passport'
import { BasicStrategy } from 'passport-http'

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
      } else {
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
    }
  )(req, res)
})

const port = process.env.PORT || 9988
server.listen(port, () => {
  console.log(`WebDAV server is listening on port ${port}`)
})

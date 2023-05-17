import { createServer, IncomingMessage, ServerResponse } from 'http'
import passport from 'passport'
import { BasicStrategy } from 'passport-http'
import handleCopy from '@/method/handle-copy'
import handleDelete from '@/method/handle-delete'
import handleGet from '@/method/handle-get'
import handleHead from '@/method/handle-head'
import handleMkcol from '@/method/handle-mkcol'
import handleMove from '@/method/handle-move'
import handleOptions from '@/method/handle-options'
import handlePropfind from '@/method/handle-propfind'
import handleProppatch from '@/method/handle-proppatch'
import handlePut from '@/method/handle-put'

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

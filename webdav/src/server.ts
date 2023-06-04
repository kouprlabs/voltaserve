import '@/infra/env'
import { createServer, IncomingMessage, ServerResponse } from 'http'
import passport from 'passport'
import { BasicStrategy } from 'passport-http'
import { Token } from '@/api/token'
import { IDP_URL, PORT } from '@/config/config'
import handleCopy from '@/handler/handle-copy'
import handleDelete from '@/handler/handle-delete'
import handleGet from '@/handler/handle-get'
import handleHead from '@/handler/handle-head'
import handleMkcol from '@/handler/handle-mkcol'
import handleMove from '@/handler/handle-move'
import handleOptions from '@/handler/handle-options'
import handlePropfind from '@/handler/handle-propfind'
import handleProppatch from '@/handler/handle-proppatch'
import handlePut from '@/handler/handle-put'

const tokens = new Map<string, Token>()

passport.use(
  new BasicStrategy(async (username, password, done) => {
    const formBody = []
    formBody.push('grant_type=password')
    formBody.push(`username=${encodeURIComponent(username)}`)
    formBody.push(`password=${encodeURIComponent(password)}`)
    try {
      const result = await fetch(`${IDP_URL}/v1/token`, {
        method: 'POST',
        body: formBody.join('&'),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      })
      const token = await result.json()
      tokens.set(username, token)
      return done(null, token)
    } catch (err) {
      return done(err, false)
    }
  })
)

const server = createServer((req: IncomingMessage, res: ServerResponse) => {
  if (req.url === '/v1/health' && req.method === 'GET') {
    res.statusCode = 200
    res.end('OK')
    return
  }
  passport.authenticate(
    'basic',
    { session: false },
    async (err: Error, token: Token) => {
      if (err || !token) {
        res.statusCode = 401
        res.setHeader('WWW-Authenticate', 'Basic realm="WebDAV Server"')
        res.end()
      } else {
        const method = req.method
        console.log(method)
        switch (method) {
          case 'OPTIONS':
            await handleOptions(req, res)
            break
          case 'GET':
            await handleGet(req, res, token)
            break
          case 'HEAD':
            await handleHead(req, res, token)
            break
          case 'PUT':
            await handlePut(req, res, token)
            break
          case 'DELETE':
            await handleDelete(req, res, token)
            break
          case 'MKCOL':
            await handleMkcol(req, res, token)
            break
          case 'COPY':
            await handleCopy(req, res, token)
            break
          case 'MOVE':
            await handleMove(req, res, token)
            break
          case 'PROPFIND':
            await handlePropfind(req, res, token)
            break
          case 'PROPPATCH':
            await handleProppatch(req, res)
            break
          default:
            res.statusCode = 501
            res.end()
        }
      }
    }
  )(req, res)
})

server.listen(PORT, () => {
  console.log(`Listening on port ${PORT}`)
})

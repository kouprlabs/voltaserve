import { createServer, IncomingMessage, ServerResponse } from 'http'
import fetch from 'node-fetch'
import passport from 'passport'
import { BasicStrategy } from 'passport-http'
import { Token } from '@/api/token'
import { IDP_URL, PORT } from '@/config/config'
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
            await handlePropfind(req, res, token)
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

server.listen(PORT, () => {
  console.log(`WebDAV server is listening on port ${PORT}`)
})

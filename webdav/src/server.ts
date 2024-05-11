import '@/infra/env'
import { createServer, IncomingMessage, ServerResponse } from 'http'
import passport from 'passport'
import { BasicStrategy } from 'passport-http'
import { TokenAPI, Token } from '@/client/idp'
import { PORT } from '@/config'
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
import { newExpiry } from '@/helper/token'

const tokens = new Map<string, Token>()
const expiries = new Map<string, Date>()
const api = new TokenAPI()

/* Refresh tokens */
setInterval(async () => {
  for (const [username, token] of tokens) {
    const expiry = expiries.get(username)
    const earlyExpiry = new Date(expiry)
    earlyExpiry.setMinutes(earlyExpiry.getMinutes() - 1)
    if (new Date() >= earlyExpiry) {
      const newToken = await api.exchange({
        grant_type: 'refresh_token',
        refresh_token: token.refresh_token,
      })
      tokens.set(username, newToken)
      expiries.set(username, newExpiry(newToken))
    }
  }
}, 5000)

passport.use(
  new BasicStrategy(async (username, password, done) => {
    try {
      let token = tokens.get(username)
      if (!token) {
        token = await new TokenAPI().exchange({
          username,
          password,
          grant_type: 'password',
        })
        tokens.set(username, token)
        expiries.set(username, newExpiry(token))
      }
      return done(null, token)
    } catch (err) {
      return done(err, false)
    }
  }),
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
        const method = req.method.toUpperCase()
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
            debugger
            res.statusCode = 501
            res.end()
        }
      }
    },
  )(req, res)
})

server.listen(PORT, () => {
  console.log(`Listening on port ${PORT}`)
})

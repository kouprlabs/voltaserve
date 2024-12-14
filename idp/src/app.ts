// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import '@/infra/env.ts'
import { Hono } from 'hono'
import { jwt } from 'hono/jwt'
import accountRouter from '@/account/router.ts'
import { getConfig } from '@/config/config.ts'
import healthRouter from '@/health/router.ts'
import {
  ErrorCode,
  ErrorData,
  newError,
  newResponse,
} from '@/infra/error/core.ts'
import tokenRouter from '@/token/router.ts'
import userRepo from '@/user/repo.ts'
import userRouter from '@/user/router.ts'
import versionRouter from '@/version/router.ts'
import { client as postgres } from '@/infra/postgres.ts'
import process from 'node:process'
import { User } from '@/user/model.ts'

const app = new Hono()

app.onError((error, c) => {
  if (error.name === 'UnauthorizedError') {
    return c.json(
      newResponse(newError({ code: ErrorCode.InvalidCredentials })),
      401,
    )
  } else {
    const genericError = error as any
    if (
      genericError.code && Object.values(ErrorCode).includes(genericError.code)
    ) {
      const data = genericError as ErrorData
      if (data.error) {
        console.error(data.error)
      }
      return c.json(newResponse(data), data.status)
    } else {
      console.error(error)
      return c.json(
        newResponse(newError({ code: ErrorCode.InternalServerError })),
        500,
      )
    }
  }
})

app.route('/version', versionRouter)
app.route('/v3/accounts', accountRouter)

app.use('/v3/*', (c, next) => {
  const jwtMiddleware = jwt({ secret: getConfig().token.jwtSigningKey })
  switch (c.req.path) {
    case '/v3/users/me/picture:extension':
      return next()
    default:
      return jwtMiddleware(c, next)
  }
})

declare module 'hono' {
  interface ContextVariableMap {
    user: User
  }
}

app.use('/v3/*', async (c, next) => {
  const payload = c.get('jwtPayload')
  try {
    const user = await userRepo.findById(payload.sub)
    c.set('user', user)
  } catch {
    return c.json(
      newResponse(newError({ code: ErrorCode.InvalidCredentials })),
      401,
    )
  } finally {
    await next()
  }
})

app.route('/v3/health', healthRouter)
app.route('/v3/users', userRouter)
app.route('/v3/token', tokenRouter)

postgres
  .connect()
  .then(() => Deno.serve({ port: getConfig().port }, app.fetch))
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })

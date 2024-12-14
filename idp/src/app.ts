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
import { logger } from 'hono/logger'
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

app.use(logger())

app.onError((error, c) => {
  if (error.message === 'Unauthorized') {
    return c.json(
      newResponse(newError({ code: ErrorCode.InvalidCredentials })),
      401,
    )
  } else {
    console.error(error)
    return c.json(
      newResponse(newError({ code: ErrorCode.InternalServerError })),
      500,
    )
  }
})

app.use('*', async (c, next) => {
  try {
    return await next()
  } catch (error: any) {
    if (error.code && Object.values(ErrorCode).includes(error.code)) {
      const data = error as ErrorData
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

app.use('/v3/*', async (c, next) => {
  const jwtMiddleware = jwt({ secret: getConfig().token.jwtSigningKey })
  if (
    c.req.path.startsWith('/v3/accounts') ||
    c.req.path.startsWith('/v3/users/me/picture') ||
    c.req.path === '/v3/token' ||
    c.req.path === '/v3/health' ||
    c.req.path === '/version'
  ) {
    return await next()
  } else {
    return await jwtMiddleware(c, next)
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
app.route('/v3/accounts', accountRouter)
app.route('/v3/token', tokenRouter)
app.route('/version', versionRouter)

postgres
  .connect()
  .then(() => Deno.serve({ port: getConfig().port }, app.fetch))
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })

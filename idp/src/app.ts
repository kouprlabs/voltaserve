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
import accountRouter from '@/account/router.ts'
import { getConfig } from '@/config/config.ts'
import healthRouter from '@/health/router.ts'
import {
  ErrorCode,
  errorHandler,
  newError,
  newResponse,
} from '@/infra/error/core.ts'
import tokenRouter from '@/token/router.ts'
import userRepo from '@/user/repo.ts'
import userRouter from '@/user/router.ts'
import versionRouter from '@/version/router.ts'
import { client as postgres } from '@/infra/postgres.ts'
import process from 'node:process'

const app = new Hono()

const { jwtSigningKey: secretOrKey, issuer, audience } = getConfig().token
passport.use(
  new JwtStrategy(
    {
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey,
      issuer,
      audience,
    },
    async (payload: any, done: any) => {
      try {
        const user = await userRepo.findById(payload.sub)
        return done(null, user)
      } catch {
        return done(
          newResponse(newError({ code: ErrorCode.InvalidCredentials })),
          false,
        )
      }
    },
  ),
)

app.route('/v3/health', healthRouter)
app.route('/v3/users', userRouter)
app.route('/v3/accounts', accountRouter)
app.route('/v3/token', tokenRouter)
app.route('/version', versionRouter)

const port = getConfig().port

postgres
  .connect()
  .then(() => Deno.serve({ port }, app.fetch))
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })

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
import cors from 'cors'
import express from 'express'
import logger from 'morgan'
import passport from 'passport'
import { ExtractJwt, Strategy as JwtStrategy } from 'passport-jwt'
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
import { client as postgres } from './infra/postgres.ts'
import process from 'node:process'

const app = express()

app.use(cors())
app.use(logger('dev'))
app.use(express.json({ limit: '3mb' }))
app.use(express.urlencoded({ extended: true }))

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

app.use('/v3/health', healthRouter)
app.use('/v3/users', userRouter)
app.use('/v3/accounts', accountRouter)
app.use('/v3/token', tokenRouter)
app.use('/version', versionRouter)

app.use(errorHandler)

const port = getConfig().port

postgres
  .connect()
  .then(() => {
    app.listen(port, () => {
      console.log(`Listening on port ${port}`)
    })
  })
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })

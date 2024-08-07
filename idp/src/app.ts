// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import '@/infra/env'
import bodyParser from 'body-parser'
import cors from 'cors'
import express from 'express'
import logger from 'morgan'
import passport from 'passport'
import { Strategy as JwtStrategy, ExtractJwt } from 'passport-jwt'
import accountRouter from '@/account/router'
import { getConfig } from '@/config/config'
import healthRouter from '@/health/router'
import { errorHandler } from '@/infra/error'
import tokenRouter from '@/token/router'
import userRepo from '@/user/repo'
import userRouter from '@/user/router'
import { client as postgres } from './infra/postgres'

const app = express()

app.use(cors())
app.use(logger('dev'))
app.use(express.json({ limit: '3mb' }))
app.use(express.urlencoded({ extended: true }))
app.use(bodyParser.json())

const { jwtSigningKey: secretOrKey, issuer, audience } = getConfig().token
passport.use(
  new JwtStrategy(
    {
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey,
      issuer,
      audience,
    },
    async (payload, done) => {
      try {
        const user = await userRepo.findByID(payload.sub)
        return done(null, user)
      } catch {
        return done(null, false)
      }
    },
  ),
)

app.use('/v2/health', healthRouter)
app.use('/v2/user', userRouter)
app.use('/v2/accounts', accountRouter)
app.use('/v2/token', tokenRouter)

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

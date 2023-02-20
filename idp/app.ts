import bodyParser from 'body-parser'
import cors from 'cors'
import dotenv from 'dotenv'
import logger from 'morgan'
import { Strategy as JwtStrategy, ExtractJwt } from 'passport-jwt'
import { URL } from 'url'
import express from 'express'
import passport from 'passport'
import accountRouter from './account/router'
import { getConfig } from './infra/config'
import { UserRepo } from './infra/db'
import { errorHandler } from './infra/error'
import tokenRouter from './token/router'
import userRouter from './user/router'

dotenv.config()

const app = express()

app.use(cors())
app.use(logger('dev'))
app.use(express.json({ limit: '3mb' }))
app.use(express.urlencoded({ extended: true }))
app.use(bodyParser.json())

const tokenConfig = getConfig().token
passport.use(
  new JwtStrategy(
    {
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey: tokenConfig.jwtSigningKey,
      issuer: tokenConfig.issuer,
      audience: tokenConfig.audience,
    },
    async (jwt_payload, done) => {
      const user = await UserRepo.find('id', jwt_payload.sub)
      if (user) {
        return done(null, user)
      } else {
        return done(null, false)
      }
    }
  )
)

app.use('/v1/user', userRouter)
app.use('/v1/accounts', accountRouter)
app.use('/v1/token', tokenRouter)

app.use(errorHandler)

const port = new URL(getConfig().url).port

app.listen(port, () => {
  console.log(`Listening on port ${port}`)
})

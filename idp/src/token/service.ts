import jwt from 'jsonwebtoken'
import { getConfig } from '@/config/config'
import { ErrorCode, newError } from '@/infra/error'
import { newHyphenlessUuid } from '@/infra/id'
import { verifyPassword } from '@/infra/password'
import userRepo from '@/user/repo'

export type TokenGrantType = 'password' | 'refresh_token'

// https://datatracker.ietf.org/doc/html/rfc6749#section-5.1
export type Token = {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token: string
}

// https://datatracker.ietf.org/doc/html/rfc6749#section-4.3.2
export type TokenExchangeOptions = {
  grant_type: TokenGrantType
  username?: string
  password?: string
  refresh_token?: string
}

export async function exchange(options: TokenExchangeOptions): Promise<Token> {
  validateParemeters(options)
  // https://datatracker.ietf.org/doc/html/rfc6749#section-4.3
  if (options.grant_type === 'password') {
    const user = await userRepo.find('username', options.username)
    if (!user) {
      throw newError({ code: ErrorCode.InvalidUsernameOrPassword })
    }
    if (!user.isEmailConfirmed) {
      throw newError({ code: ErrorCode.EmailNotConfimed })
    }
    if (verifyPassword(options.password, user.passwordHash)) {
      return newToken(user.id)
    } else {
      throw newError({ code: ErrorCode.InvalidUsernameOrPassword })
    }
  }
  // https://datatracker.ietf.org/doc/html/rfc6749#section-6
  if (options.grant_type === 'refresh_token') {
    const user = await userRepo.find(
      'refresh_token_value',
      options.refresh_token
    )
    if (!user) {
      throw newError({ code: ErrorCode.InvalidUsernameOrPassword })
    }
    if (!user.isEmailConfirmed) {
      throw newError({ code: ErrorCode.EmailNotConfimed })
    }
    return newToken(user.id)
  }
}

function validateParemeters(options: TokenExchangeOptions) {
  if (!options.grant_type) {
    throw newError({
      code: ErrorCode.InvalidRequest,
      message: 'Missing parameter: grant_type',
    })
  }
  if (
    options.grant_type !== 'password' &&
    options.grant_type !== 'refresh_token'
  ) {
    throw newError({
      code: ErrorCode.UnsupportedGrantType,
      message: `Grant type unsupported: ${options.grant_type}`,
    })
  }
  if (options.grant_type === 'password') {
    if (!options.username) {
      throw newError({
        code: ErrorCode.InvalidRequest,
        message: 'Missing parameter: username',
      })
    }
    if (!options.password) {
      throw newError({
        code: ErrorCode.InvalidRequest,
        message: 'Missing parameter: password',
      })
    }
  }
  if (options.grant_type === 'refresh_token' && !options.refresh_token) {
    throw newError({
      code: ErrorCode.InvalidRequest,
      message: 'Missing parameter: refresh_token',
    })
  }
}

async function newToken(userId: string): Promise<Token> {
  const config = getConfig().token
  const token: Token = {
    access_token: jwt.sign({}, config.jwtSigningKey, {
      expiresIn: newAccessTokenExpiry(),
      audience: config.audience,
      issuer: config.issuer,
      subject: userId,
    }),
    expires_in: config.accessTokenLifetime,
    token_type: 'Bearer',
    refresh_token: newHyphenlessUuid(),
  }
  const user = await userRepo.find('id', userId, true)
  await userRepo.update({
    id: user.id,
    refreshTokenValue: token.refresh_token,
    refreshTokenValidTo: newRefreshTokenExpiry(),
  })
  return token
}

function newAccessTokenExpiry(): number {
  const now = new Date()
  now.setSeconds(now.getSeconds() + getConfig().token.accessTokenLifetime)
  return Math.floor(now.getTime() / 1000)
}

function newRefreshTokenExpiry(): number {
  const now = new Date()
  now.setSeconds(now.getSeconds() + getConfig().token.refreshTokenLifetime)
  return Math.floor(now.getTime() / 1000)
}

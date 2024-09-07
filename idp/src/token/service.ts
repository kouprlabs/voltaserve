// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { decodeJwt, SignJWT } from 'jose'
import { getConfig } from '@/config/config'
import { ErrorCode, newError } from '@/infra/error'
import { newHyphenlessUuid } from '@/infra/id'
import { verifyPassword } from '@/infra/password'
import { User } from '@/user/model'
import userRepo from '@/user/repo'
import {getUserByAdmin} from "@/user/service";

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
  validateParameters(options)
  // https://datatracker.ietf.org/doc/html/rfc6749#section-4.3
  if (options.grant_type === 'password') {
    let user: User
    try {
      user = await userRepo.findByUsername(options.username.toLocaleLowerCase())
    } catch {
      throw newError({ code: ErrorCode.InvalidUsernameOrPassword })
    }
    if (!user.isEmailConfirmed) {
      throw newError({ code: ErrorCode.EmailNotConfimed })
    }
    if (!user.isActive) {
      throw newError({ code: ErrorCode.UserSuspended })
    }
    if (verifyPassword(options.password, user.passwordHash)) {
      return newToken(user.id, user.isAdmin)
    } else {
      throw newError({ code: ErrorCode.InvalidUsernameOrPassword })
    }
  }
  // https://datatracker.ietf.org/doc/html/rfc6749#section-6
  if (options.grant_type === 'refresh_token') {
    let user: User
    try {
      user = await userRepo.findByRefreshTokenValue(options.refresh_token)
    } catch {
      throw newError({ code: ErrorCode.InvalidUsernameOrPassword })
    }
    if (!user.isEmailConfirmed) {
      throw newError({ code: ErrorCode.EmailNotConfimed })
    }
    if (new Date() >= new Date(user.refreshTokenExpiry)) {
      throw newError({ code: ErrorCode.RefreshTokenExpired })
    }
    return newToken(user.id, user.isAdmin)
  }
}

export const checkAdmin = (jwt) => {
  if (!decodeJwt(jwt).is_admin)
    throw newError({ code: ErrorCode.MissingPermission })
}

export const checkForcePasswordChange = async (userId: string) => {
  const user = await getUserByAdmin(userId)
  if (user.forceChangePassword) {
    const resetPasswordToken = await userRepo.getResetPasswordToken(user.id)
    throw newError({
      code: ErrorCode.ForceChangePassword,
      message: resetPasswordToken,
    })
  }
}

function validateParameters(options: TokenExchangeOptions) {
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

async function newToken(userId: string, isAdmin: boolean): Promise<Token> {
  const config = getConfig().token
  const expiry = newAccessTokenExpiry()
  const jwt = await new SignJWT({ sub: userId, is_admin: isAdmin })
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setIssuer(config.issuer)
    .setAudience(config.audience)
    .setExpirationTime(expiry)
    .sign(new TextEncoder().encode(config.jwtSigningKey))
  const token: Token = {
    access_token: jwt,
    expires_in: expiry,
    token_type: 'Bearer',
    refresh_token: newHyphenlessUuid(),
  }
  const user = await userRepo.findByID(userId)
  await userRepo.update({
    id: user.id,
    refreshTokenValue: token.refresh_token,
    refreshTokenExpiry: newRefreshTokenExpiry(),
  })
  return token
}

function newRefreshTokenExpiry(): string {
  const now = new Date()
  now.setSeconds(now.getSeconds() + getConfig().token.refreshTokenLifetime)
  return now.toISOString()
}

function newAccessTokenExpiry(): number {
  const now = new Date()
  now.setSeconds(now.getSeconds() + getConfig().token.refreshTokenLifetime)
  return Math.floor(now.getTime() / 1000)
}

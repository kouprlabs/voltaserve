// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { decode, sign, verify } from 'hono/jwt'
import { getConfig } from '@/config/config.ts'
import {
  newEmailNotConfirmedError,
  newInvalidAppleTokenError,
  newInvalidGrantType,
  newInvalidUsernameOrPasswordError,
  newMissingFormParamError,
  newRefreshTokenExpiredError,
  newUserNotFoundError,
  newUserSuspendedError,
  newUserTemporarilyLockedError,
} from '@/error/creators.ts'
import { newHyphenlessUuid } from '@/infra/id.ts'
import { verifyPassword } from '@/infra/password.ts'
import { getApplePublicKey } from '@/infra/apple.ts'
import { User } from '@/user/model.ts'
import userRepo from '@/user/repo.ts'
import { signUpWithApple } from '@/account/service.ts'

export type TokenGrantType = 'password' | 'refresh_token' | 'apple'

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
  apple_jwt?: string
  apple_full_name?: string
}

export async function exchange(options: TokenExchangeOptions): Promise<Token> {
  validateParameters(options)
  if (options.grant_type === 'password') {
    return await exchangeWithPassword(options)
  } else if (options.grant_type === 'refresh_token') {
    return await exchangeWithRefreshToken(options)
  } else if (options.grant_type == 'apple') {
    return await exchangeWithApple(options)
  } else {
    throw newInvalidGrantType(options.grant_type)
  }
}

// https://datatracker.ietf.org/doc/html/rfc6749#section-4.3
async function exchangeWithPassword(
  options: TokenExchangeOptions,
): Promise<Token> {
  let user: User
  try {
    user = await userRepo.findByUsername(options.username!.toLocaleLowerCase())
  } catch {
    throw newInvalidUsernameOrPasswordError()
  }
  if (!user.isEmailConfirmed) {
    throw newEmailNotConfirmedError()
  }
  if (!user.isActive) {
    throw newUserSuspendedError()
  }
  if (isStillLocked(user)) {
    throw newUserTemporarilyLockedError()
  } else {
    if (verifyPassword(options.password!, user.passwordHash!)) {
      await resetFailedAttemptsAndUnlock(user.id)
      return newToken(user.id, user.isAdmin)
    } else {
      await increaseFailedAttemptsOrLock(user.id)
      throw newInvalidUsernameOrPasswordError()
    }
  }
}

// https://datatracker.ietf.org/doc/html/rfc6749#section-6
async function exchangeWithRefreshToken(
  options: TokenExchangeOptions,
): Promise<Token> {
  let user: User
  try {
    user = await userRepo.findByRefreshTokenValue(options.refresh_token!)
  } catch {
    throw newInvalidUsernameOrPasswordError()
  }
  if (!user.isEmailConfirmed) {
    throw newEmailNotConfirmedError()
  }
  if (new Date() >= new Date(user.refreshTokenExpiry!)) {
    throw newRefreshTokenExpiredError()
  }
  return newToken(user.id, user.isAdmin)
}

async function exchangeWithApple(
  options: TokenExchangeOptions,
): Promise<Token> {
  const { apple_jwt } = options
  const { header } = decode(apple_jwt!)
  const appleKey = await getApplePublicKey(header)
  let payload: any
  try {
    payload = await verify(apple_jwt!, appleKey, 'RS256')
  } catch {
    throw newInvalidAppleTokenError()
  }
  let user: User
  try {
    user = await userRepo.findByUsername(payload.sub)
  } catch (error: any) {
    if (error.code === newUserNotFoundError().code) {
      user = await signUpWithApple({
        payload,
        appleFullName: options.apple_full_name,
      })
    } else {
      throw error
    }
  }
  return newToken(user.id, user.isAdmin)
}

function validateParameters(options: TokenExchangeOptions) {
  if (!options.grant_type) {
    throw newMissingFormParamError('grant_type')
  }
  if (
    options.grant_type !== 'password' &&
    options.grant_type !== 'refresh_token' &&
    options.grant_type !== 'apple'
  ) {
    throw newInvalidGrantType(options.grant_type)
  }
  if (options.grant_type === 'password') {
    if (!options.username) {
      throw newMissingFormParamError('username')
    }
    if (!options.password) {
      throw newMissingFormParamError('password')
    }
  }
  if (options.grant_type === 'refresh_token' && !options.refresh_token) {
    throw newMissingFormParamError('refresh_token')
  }
}

async function newToken(userId: string, isAdmin: boolean): Promise<Token> {
  const config = getConfig().token
  const expiry = newAccessTokenExpiry()
  const jwt = await sign(
    {
      sub: userId,
      is_admin: isAdmin,
      exp: expiry,
      aud: config.audience,
      iss: config.issuer,
      iat: Math.floor(new Date().getTime() / 1000),
    },
    config.jwtSigningKey,
    'HS256',
  )
  const token: Token = {
    access_token: jwt,
    expires_in: expiry,
    token_type: 'Bearer',
    refresh_token: newHyphenlessUuid(),
  }
  const user = await userRepo.findById(userId)
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
  now.setSeconds(now.getSeconds() + getConfig().token.accessTokenLifetime)
  return Math.floor(now.getTime() / 1000)
}

async function increaseFailedAttemptsOrLock(userId: string): Promise<void> {
  const user = await userRepo.findById(userId)
  const failedAttempts = user.failedAttempts + 1
  if (failedAttempts <= getConfig().security.maxFailedAttempts) {
    await userRepo.update({
      id: user.id,
      failedAttempts,
    })
  } else {
    await userRepo.update({
      id: user.id,
      lockedUntil: newLockoutUntil(),
    })
  }
}

async function resetFailedAttemptsAndUnlock(userId: string): Promise<void> {
  const user = await userRepo.findById(userId)
  await userRepo.update({
    id: user.id,
    failedAttempts: 0,
    lockedUntil: null,
  })
}

function newLockoutUntil(): string {
  const now = new Date()
  now.setSeconds(now.getSeconds() + getConfig().security.lockoutPeriod)
  return now.toISOString()
}

function isStillLocked(user: User): boolean {
  return !!(user.lockedUntil && new Date() < new Date(user.lockedUntil))
}

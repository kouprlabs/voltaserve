// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Request, Response, NextFunction } from 'express'

export enum ErrorCode {
  InternalServerError = 'internal_server_error',
  RequestValidationError = 'request_validation_error',
  UsernameUnavailable = 'username_unavailable',
  ResourceNotFound = 'resource_not_found',
  InvalidUsernameOrPassword = 'invalid_username_or_password',
  InvalidPassword = 'invalid_password',
  InvalidJwt = 'invalid_jwt',
  EmailNotConfimed = 'email_not_confirmed',
  RefreshTokenExpired = 'refresh_token_expired',
  InvalidRequest = 'invalid_request',
  UnsupportedGrantType = 'unsupported_grant_type',
  PasswordValidationFailed = 'password_validation_failed',
  UserSuspended = 'user_suspended',
  MissingPermission = 'missing_permission',
}

const statuses: { [key: string]: number } = {
  [ErrorCode.InternalServerError]: 500,
  [ErrorCode.RequestValidationError]: 400,
  [ErrorCode.UsernameUnavailable]: 409,
  [ErrorCode.ResourceNotFound]: 404,
  [ErrorCode.InvalidUsernameOrPassword]: 401,
  [ErrorCode.InvalidPassword]: 401,
  [ErrorCode.InvalidJwt]: 401,
  [ErrorCode.EmailNotConfimed]: 401,
  [ErrorCode.RefreshTokenExpired]: 401,
  [ErrorCode.InvalidRequest]: 400,
  [ErrorCode.UnsupportedGrantType]: 400,
  [ErrorCode.PasswordValidationFailed]: 400,
  [ErrorCode.UserSuspended]: 403,
  [ErrorCode.MissingPermission]: 403,
}

const userMessages: { [key: string]: string } = {
  [ErrorCode.UsernameUnavailable]: 'Email belongs to an existing user.',
  [ErrorCode.EmailNotConfimed]: 'Email not confirmed.',
  [ErrorCode.InvalidPassword]: 'Invalid password.',
  [ErrorCode.InvalidUsernameOrPassword]: 'Invalid username or password.',
  [ErrorCode.UserSuspended]: 'User suspended.',
  [ErrorCode.MissingPermission]: 'You are not an admin',
}

export type ErrorData = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  error?: any
}

export type ErrorResponse = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
}

export type ErrorOptions = {
  code: ErrorCode
  message?: string
  userMessage?: string
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  error?: any
}

export function newError(options: ErrorOptions): ErrorData {
  const userMessage =
    options.userMessage ||
    userMessages[options.code] ||
    'Oops! something went wrong'
  return {
    code: options.code,
    status: statuses[options.code],
    message: options.message || userMessage,
    userMessage,
    moreInfo: `https://voltaserve.com/docs/idp/errors/${options.code}`,
    error: options.error,
  }
}

export function newResponse(data: ErrorData): ErrorResponse {
  return {
    code: data.code,
    status: data.status,
    message: data.message,
    userMessage: data.userMessage,
    moreInfo: data.moreInfo,
  }
}

export function errorHandler(
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  error: any,
  _: Request,
  res: Response,
  next: NextFunction,
) {
  if (error.code && Object.values(ErrorCode).includes(error.code)) {
    const data = error as ErrorData
    if (data.error) {
      console.error(data.error)
    }
    res.status(data.status).json(newResponse(data))
  } else {
    console.error(error)
    res
      .status(500)
      .json(newResponse(newError({ code: ErrorCode.InternalServerError })))
  }
  next(error)
  return
}

/* eslint-disable-next-line @typescript-eslint/no-explicit-any */
export function parseValidationError(result: any): ErrorData {
  let message: string
  let userMessage: string
  if (result.errors) {
    message = result.errors
      /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
      .map((e: any) => `${e.msg} for ${e.type} ${e.path} in ${e.location}.`)
      .join(' ')
    userMessage = result.errors
      /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
      .map((e: any) => `${e.msg} for ${e.type} ${e.path}.`)
      .join(' ')
  }
  return newError({
    code: ErrorCode.RequestValidationError,
    message,
    userMessage,
  })
}

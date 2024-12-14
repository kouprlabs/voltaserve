// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

export enum ErrorCode {
  InternalServerError = 'internal_server_error',
  RequestValidationError = 'request_validation_error',
  UsernameUnavailable = 'username_unavailable',
  ResourceNotFound = 'resource_not_found',
  InvalidUsernameOrPassword = 'invalid_username_or_password',
  InvalidPassword = 'invalid_password',
  InvalidJwt = 'invalid_jwt',
  InvalidCredentials = 'invalid_credentials',
  EmailNotConfirmed = 'email_not_confirmed',
  RefreshTokenExpired = 'refresh_token_expired',
  InvalidRequest = 'invalid_request',
  InvalidGrantType = 'invalid_grant_type',
  PasswordValidationFailed = 'password_validation_failed',
  UserSuspended = 'user_suspended',
  UserTemporarilyLocked = 'user_locked',
  UserIsNotAdmin = 'user_is_not_admin',
  UserNotFound = 'user_not_found',
  CannotSuspendLastAdmin = 'cannot_suspend_last_admin',
  SearchError = 'search_error',
}

type StatusCode = 500 | 400 | 401 | 403 | 404 | 409 | 429

const statuses: { [key: string]: StatusCode } = {
  [ErrorCode.InternalServerError]: 500,
  [ErrorCode.RequestValidationError]: 400,
  [ErrorCode.UsernameUnavailable]: 409,
  [ErrorCode.ResourceNotFound]: 404,
  [ErrorCode.InvalidUsernameOrPassword]: 401,
  [ErrorCode.InvalidPassword]: 401,
  [ErrorCode.InvalidJwt]: 401,
  [ErrorCode.InvalidCredentials]: 401,
  [ErrorCode.EmailNotConfirmed]: 401,
  [ErrorCode.RefreshTokenExpired]: 401,
  [ErrorCode.InvalidRequest]: 400,
  [ErrorCode.InvalidGrantType]: 400,
  [ErrorCode.PasswordValidationFailed]: 400,
  [ErrorCode.UserSuspended]: 403,
  [ErrorCode.UserTemporarilyLocked]: 429,
  [ErrorCode.UserIsNotAdmin]: 403,
  [ErrorCode.UserNotFound]: 404,
  [ErrorCode.CannotSuspendLastAdmin]: 400,
  [ErrorCode.SearchError]: 500,
}

export type ErrorData = {
  code: string
  status: StatusCode
  message: string
  userMessage: string
  moreInfo: string
  error?: any
}

export type ErrorResponse = {
  code: string
  status: StatusCode
  message: string
  userMessage: string
  moreInfo: string
}

export type ErrorOptions = {
  code: ErrorCode
  message?: string
  userMessage?: string
  error?: any
}

export function newError(options: ErrorOptions): ErrorData {
  const userMessage = options.userMessage || 'Oops! something went wrong.'
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

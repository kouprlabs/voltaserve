// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { ErrorCode, newError } from './core.ts'

export function newInternalServerError(error?: unknown) {
  return newError({
    code: ErrorCode.InternalServerError,
    error,
    message: 'Internal server error.',
    userMessage: 'Internal server error.',
  })
}

export function newUserNotFoundError() {
  return newError({
    code: ErrorCode.ResourceNotFound,
    message: 'User not found.',
    userMessage: 'User not found.',
  })
}

export function newPictureNotFoundError() {
  return newError({
    code: ErrorCode.ResourceNotFound,
    message: 'Picture not found.',
    userMessage: 'Picture not found.',
  })
}

export function newInvalidJwtError() {
  return newError({
    code: ErrorCode.InvalidJwt,
    message: 'Invalid JWT.',
    userMessage: 'Invalid JWT.',
  })
}

export function newUsernameUnavailableError() {
  return newError({
    code: ErrorCode.UsernameUnavailable,
    message: `Username not available.`,
    userMessage: `Username not available.`,
  })
}

export function newInvalidUsernameOrPasswordError() {
  return newError({
    code: ErrorCode.InvalidUsernameOrPassword,
    message: 'Invalid username or password.',
    userMessage: 'Invalid username or password.',
  })
}

export function newEmailNotConfirmedError() {
  return newError({
    code: ErrorCode.EmailNotConfirmed,
    message: 'Email not confirmed.',
    userMessage: 'Email not confirmed.',
  })
}

export function newUserSuspendedError() {
  return newError({
    code: ErrorCode.UserSuspended,
    message: 'User suspended.',
    userMessage: 'User suspended.',
  })
}

export function newUserTemporarilyLockedError() {
  return newError({
    code: ErrorCode.UserTemporarilyLocked,
    message: 'User temporarily locked. Try again later.',
    userMessage: 'User temporarily locked. Try again later.',
  })
}

export function newRefreshTokenExpiredError() {
  return newError({
    code: ErrorCode.RefreshTokenExpired,
    message: 'Refresh token expired.',
    userMessage: 'Refresh token expired.',
  })
}

export function newUserIsNotAdminError() {
  return newError({
    code: ErrorCode.UserIsNotAdmin,
    message: 'User is not admin.',
    userMessage: 'User is not admin.',
  })
}

export function newMissingQueryParamError(param: string) {
  return newError({
    code: ErrorCode.InvalidRequest,
    message: `Missing query param: ${param}.`,
    userMessage: `Missing query param: ${param}.`,
  })
}

export function newMissingFormParamError(param: string) {
  return newError({
    code: ErrorCode.InvalidRequest,
    message: `Missing form parameter: ${param}.`,
    userMessage: `Missing form parameter: ${param}.`,
  })
}

export function newInvalidGrantType(grantType: string) {
  return newError({
    code: ErrorCode.InvalidGrantType,
    message: `Invalid Grant type: ${grantType}.`,
    userMessage: `Invalid Grant type: ${grantType}.`,
  })
}

export function newPasswordValidationFailedError() {
  return newError({
    code: ErrorCode.PasswordValidationFailed,
    message: 'Password validation failed.',
    userMessage: 'Password validation failed.',
  })
}

export function newCannotSuspendSoleAdminError() {
  return newError({
    code: ErrorCode.CannotSuspendLastAdmin,
    message: 'Cannot suspend sole admin.',
    userMessage: 'Cannot suspend sole admin.',
  })
}

export function newCannotDemoteSoleAdminError() {
  return newError({
    code: ErrorCode.CannotSuspendLastAdmin,
    message: 'Cannot demote sole admin.',
    userMessage: 'Cannot demote sole admin.',
  })
}

export function newInvalidAppleTokenError() {
  return newError({
    code: ErrorCode.InvalidJwt,
    message: 'Invalid Apple identity token.',
    userMessage: 'Unable to verify your Apple sign-in. Please try again.',
  })
}

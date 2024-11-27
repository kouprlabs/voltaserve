import { ErrorCode, newError } from './core'

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

export function newInvalidPasswordError() {
  return newError({
    code: ErrorCode.InvalidPassword,
    message: 'Invalid password.',
    userMessage: 'Invalid password.',
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

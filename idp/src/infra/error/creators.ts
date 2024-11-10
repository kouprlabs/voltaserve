import { ErrorCode, newError } from './core'

export function newUserNotFoundError(field: string, value: string) {
  return newError({
    code: ErrorCode.ResourceNotFound,
    message: `User with ${field} '${value}' not found.`,
    userMessage: `User with ${field} '${value}' not found.`,
  })
}

export function newUserInsertError(username: string) {
  return newError({
    code: ErrorCode.InternalServerError,
    message: `Failed to insert user with username '${username}'.`,
    userMessage: `Failed to insert user with username '${username}'.`,
  })
}

export function newUserUpdateError(id: string) {
  return newError({
    code: ErrorCode.InternalServerError,
    message: `Failed to update user with id '${id}'.`,
    userMessage: `Failed to update user with id '${id}'.`,
  })
}

export function newUserPictureNotFoundError(id: string) {
  return newError({
    code: ErrorCode.ResourceNotFound,
    message: `User picture for id '${id}' not found.`,
    userMessage: `User picture for id '${id}' not found.`,
  })
}

export function newInvalidJwtError() {
  return newError({
    code: ErrorCode.InvalidJwt,
    message: 'Invalid JWT.',
    userMessage: 'Invalid JWT.',
  })
}

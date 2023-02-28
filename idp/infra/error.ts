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
  InvalidRequest = 'invalid_request',
  UnsupportedGrantType = 'unsupported_grant_type',
  PasswordValidationFailed = 'password_validation_failed',
}

const statusMap: { [key: string]: number } = {}
statusMap[ErrorCode.InternalServerError] = 500
statusMap[ErrorCode.RequestValidationError] = 400
statusMap[ErrorCode.UsernameUnavailable] = 409
statusMap[ErrorCode.ResourceNotFound] = 404
statusMap[ErrorCode.InvalidUsernameOrPassword] = 401
statusMap[ErrorCode.InvalidPassword] = 401
statusMap[ErrorCode.InvalidJwt] = 401
statusMap[ErrorCode.EmailNotConfimed] = 401
statusMap[ErrorCode.InvalidRequest] = 400
statusMap[ErrorCode.UnsupportedGrantType] = 400
statusMap[ErrorCode.PasswordValidationFailed] = 400

const userMessageMap: { [key: string]: string } = {}
userMessageMap[ErrorCode.UsernameUnavailable] =
  'Email belongs to an existing user.'
userMessageMap[ErrorCode.EmailNotConfimed] = 'Email not confirmed.'
userMessageMap[ErrorCode.InvalidPassword] = 'Invalid password.'
userMessageMap[ErrorCode.InvalidUsernameOrPassword] =
  'Invalid username or password.'

export type IdpError = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
  error?: any
}

export type NewErrorOptions = {
  code: ErrorCode
  message?: string
  userMessage?: string
  error?: any
}

export function newError(opts: NewErrorOptions): IdpError {
  return {
    code: opts.code,
    status: statusMap[opts.code] || 500,
    message: opts.message || 'Internal server error',
    userMessage:
      opts.userMessage ||
      userMessageMap[opts.code] ||
      'Oops! something went wrong',
    moreInfo: `https://voltaserve.com/docs/idp/errors/${opts.code}`,
    error: opts.error || undefined,
  }
}

export function errorHandler(
  error: any,
  _: Request,
  res: Response,
  next: NextFunction
) {
  if (error.code && Object.values(ErrorCode).includes(error.code)) {
    const e = error as IdpError
    if (e.error) {
      console.error(e.error)
    }
    res.status(e.status).json({
      code: e.code,
      message: e.message,
      userMessage: e.userMessage,
      moreInfo: e.moreInfo,
    })
  } else {
    console.error(error)
    res.status(500).json({
      code: ErrorCode.InternalServerError,
    })
  }
  next(error)
  return
}

export function parseValidationError(result: any): IdpError {
  const message = result.errors
    ? result.errors
        .map(
          (e: any) => `${e.msg} for parameter '${e.param}' in ${e.location}.`
        )
        .join(' ')
    : undefined
  return newError({ code: ErrorCode.RequestValidationError, message })
}

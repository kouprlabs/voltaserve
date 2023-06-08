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
}

const userMessages: { [key: string]: string } = {
  [ErrorCode.UsernameUnavailable]: 'Email belongs to an existing user.',
  [ErrorCode.EmailNotConfimed]: 'Email not confirmed.',
  [ErrorCode.InvalidPassword]: 'Invalid password.',
  [ErrorCode.InvalidUsernameOrPassword]: 'Invalid username or password.',
}

export type ErrorData = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
  error?: any
}

export type ErrorResponse = {
  code: string
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
    message: data.message,
    userMessage: data.userMessage,
    moreInfo: data.moreInfo,
  }
}

export function errorHandler(
  error: any,
  _: Request,
  res: Response,
  next: NextFunction
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

export function parseValidationError(result: any): ErrorData {
  let message: string
  if (result.errors) {
    message = result.errors
      .map((e: any) => `${e.msg} for parameter '${e.param}' in ${e.location}.`)
      .join(' ')
  }
  return newError({ code: ErrorCode.RequestValidationError, message })
}

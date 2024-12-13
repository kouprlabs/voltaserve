import { ErrorCode, ErrorData, newError } from '@/infra/error/core.ts'

/* eslint-disable-next-line @typescript-eslint/no-explicit-any */
export function parseValidationError(result: any): ErrorData {
  let message: string | undefined
  let userMessage: string | undefined
  if (result.errors) {
    message = result.errors
      /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
      .map((e: any) => `${e.msg} for ${e.type} '${e.path}' in ${e.location}.`)
      .join(' ')
    userMessage = result.errors
      /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
      .map((e: any) => `${e.msg} for ${e.type} '${e.path}'.`)
      .join(' ')
  }
  return newError({
    code: ErrorCode.RequestValidationError,
    message,
    userMessage,
  })
}

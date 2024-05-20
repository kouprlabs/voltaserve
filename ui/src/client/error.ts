export type ErrorResponse = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function errorToString(value: any): string {
  if (value.code && value.message && value.userMessage && value.moreInfo) {
    const error = value as ErrorResponse
    return error.userMessage
  }
  return value.toString()
}

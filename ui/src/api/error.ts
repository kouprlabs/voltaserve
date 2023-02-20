type ErrorResponse = {
  code: string
  message: string
  userMessage: string
  moreInfo: string
}

export function errorToString(value: any): string {
  if (value.code && value.message && value.userMessage && value.moreInfo) {
    const error = value as ErrorResponse
    return error.userMessage
  }
  return value.toString()
}

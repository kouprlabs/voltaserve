export type ErrorResponse = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
}

export class ClientError extends Error {
  constructor(readonly error: ErrorResponse) {
    super(error.code)
  }
}
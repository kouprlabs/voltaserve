import { ServerResponse } from 'http'
import { APIError } from '@/client/api'
import { IdPError } from '@/client/idp'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function handleError(err: any, res: ServerResponse) {
  if (err instanceof APIError) {
    res.statusCode = err.error.status
    res.statusMessage = err.error.userMessage
    res.end()
  } else if (err instanceof IdPError) {
    res.statusCode = err.error.status
    res.statusMessage = err.error.userMessage
    res.end()
  } else {
    res.statusCode = 500
    res.end()
  }
  console.error(err)
}

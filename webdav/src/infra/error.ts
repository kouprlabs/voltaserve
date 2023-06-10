import { ClientError } from "@/client/error"
import { ServerResponse } from "http"

export function handleException(err: any, res: ServerResponse) {
  if (err instanceof ClientError) {
    console.error(JSON.stringify(err.error, null, 2))
    res.statusCode = err.error.status
    res.statusMessage = err.error.userMessage
    res.end()
  } else {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}
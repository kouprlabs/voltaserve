import { IncomingMessage } from 'http'
import path from 'path'

export function getTargetPath(req: IncomingMessage) {
  const destinationHeader = req.headers.destination as string
  if (!destinationHeader) {
    return null
  }
  // Check if the destination header is a full URL
  if (
    destinationHeader.startsWith('http://') ||
    destinationHeader.startsWith('https://')
  ) {
    const url = new URL(destinationHeader)
    return path.join(decodeURIComponent(url.pathname))
  } else {
    /* Extract the path from the destination header */
    const startIndex =
      destinationHeader.indexOf(req.headers.host) + req.headers.host.length
    const value = destinationHeader.substring(startIndex)
    return path.join(decodeURIComponent(value))
  }
}

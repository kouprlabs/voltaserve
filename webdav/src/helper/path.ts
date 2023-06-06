import { IncomingMessage } from 'http'

export function getTargetPath(req: IncomingMessage) {
  const destination = req.headers.destination as string
  if (!destination) {
    return null
  }
  // Check if the destination header is a full URL
  if (destination.startsWith('http://') || destination.startsWith('https://')) {
    return new URL(destination).pathname
  } else {
    /* Extract the path from the destination header */
    const startIndex =
      destination.indexOf(req.headers.host) + req.headers.host.length
    return destination.substring(startIndex)
  }
}

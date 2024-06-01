import { encode, decode } from 'js-base64'
import { Query as FileQuery } from '@/client/api/file'

export function encodeQuery(value: string) {
  return encode(value, true)
}

export function decodeQuery(value: string): string | undefined {
  if (!value) {
    return undefined
  }
  return decode(value)
}

export function encodeFileQuery(value: FileQuery) {
  return encode(JSON.stringify(value), true)
}

export function decodeFileQuery(value: string): FileQuery | undefined {
  if (!value) {
    return undefined
  }
  return JSON.parse(decode(value))
}

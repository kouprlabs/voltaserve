import { encode, decode } from 'js-base64'

export function encodeQuery(value: string) {
  return encode(value, true)
}

export function decodeQuery(value: string): string | null {
  if (!value) {
    return null
  }
  return decode(value)
}

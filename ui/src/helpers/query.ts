import { encode, decode } from 'js-base64'

export function encodeQuery(value: string) {
  return encode(value, true)
}

export function decodeQuery(value: string): string | undefined {
  if (!value) {
    return undefined
  }
  return decode(value)
}

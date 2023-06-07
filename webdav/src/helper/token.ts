import { Token } from '@/client/idp'

export function newExpiry(token: Token): Date {
  const now = new Date()
  now.setSeconds(now.getSeconds() + token.expires_in)
  return now
}

import { scryptSync, randomBytes } from 'crypto'

export function hashPassword(password: string): string {
  const salt = randomBytes(16).toString('hex')
  const key = scryptSync(password, salt, 64).toString('hex')
  return `${key}:${salt}`
}

export function verifyPassword(password: string, hash: string): boolean {
  const [key, salt] = hash.split(':')
  const newKey = scryptSync(password, salt, 64).toString('hex')
  return newKey === key
}

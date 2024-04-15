import { v4 as uuidv4 } from 'uuid'
import hashids from 'hashids'

export function newHashId(): string {
  return new hashids(uuidv4()).encode(Date.now())
}

export function newHyphenlessUuid(): string {
  return uuidv4().replaceAll('-', '')
}

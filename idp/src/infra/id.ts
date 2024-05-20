import hashids from 'hashids'
import { v4 as uuidv4 } from 'uuid'

export function newHashId(): string {
  return new hashids(uuidv4()).encode(Date.now())
}

export function newHyphenlessUuid(): string {
  return uuidv4().replaceAll('-', '')
}

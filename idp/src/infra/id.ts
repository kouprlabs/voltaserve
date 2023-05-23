import { v4 as uuidv4 } from 'uuid'
import hashids from 'hashids/cjs'

export function newHashId(): string {
  return new hashids(uuidv4()).encode(Date.now())
}

export function newHyphenlessUuid(): string {
  return uuidv4().replaceAll('-', '')
}

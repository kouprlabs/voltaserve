import { Request } from 'express'
import { UserEntity } from './db'

export interface PassportRequest extends Request {
  user: UserEntity
}

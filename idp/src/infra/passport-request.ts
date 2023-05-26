import { Request } from 'express'
import { UserEntity } from './postgres'

export interface PassportRequest extends Request {
  user: UserEntity
}

import { Request } from 'express'
import { User } from '@/user/model'

export interface PassportRequest extends Request {
  user: User
}

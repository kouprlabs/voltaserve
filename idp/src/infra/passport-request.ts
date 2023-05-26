import { Request } from 'express'
import { User } from '@/user/repo'

export interface PassportRequest extends Request {
  user: User
}

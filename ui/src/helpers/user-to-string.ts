import { User } from '@/client/idp/user'

export default function userToString(user: User) {
  return `${user.fullName} (${user.email})`
}

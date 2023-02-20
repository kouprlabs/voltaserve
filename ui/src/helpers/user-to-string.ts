import { User } from '@/api/user'

export default function userToString(user: User) {
  return `${user.fullName} (${user.email})`
}

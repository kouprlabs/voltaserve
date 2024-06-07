import { Avatar } from '@chakra-ui/react'
import cx from 'classnames'
import { User } from '@/client/idp/user'

export type AccountMenuAvatarImageProps = {
  user: User
}

const AccountMenuAvatarImage = ({ user }: AccountMenuAvatarImageProps) => (
  <Avatar
    name={user.fullName}
    src={user.picture}
    size="sm"
    className={cx(
      'w-[40px]',
      'h-[40px]',
      'border',
      'border-gray-300',
      'dark:border-gray-700',
    )}
  />
)

export default AccountMenuAvatarImage

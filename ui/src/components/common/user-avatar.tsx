import { Avatar, SkeletonCircle } from '@chakra-ui/react'
import cx from 'classnames'
import { User } from '@/client/idp/user'

export type UserAvatarProps = {
  user?: User
  size: string
}

const UserAvatar = ({ user, size }: UserAvatarProps) => {
  if (user) {
    return (
      <Avatar
        src={user.picture}
        style={{ width: size, height: size }}
        className={cx(
          'border',
          'border-solid',
          'border-gray-300',
          'dark:border-gray-700',
        )}
      />
    )
  } else {
    return <SkeletonCircle size={size} />
  }
}

export default UserAvatar

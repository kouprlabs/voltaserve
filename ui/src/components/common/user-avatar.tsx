import {
  Avatar,
  SkeletonCircle,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { User } from '@/client/idp/user'

export type UserAvatarProps = {
  user?: User
  size: string
}

const UserAvatar = ({ user, size }: UserAvatarProps) => {
  const borderColor = useToken(
    'colors',
    useColorModeValue('gray.300', 'gray.700'),
  )
  if (user) {
    return (
      <Avatar
        src={user.picture}
        style={{ width: size, height: size }}
        border={`1px solid ${borderColor}`}
      />
    )
  } else {
    return <SkeletonCircle size={size} />
  }
}

export default UserAvatar

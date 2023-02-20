import {
  Avatar,
  SkeletonCircle,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { User } from '@/api/user'

type UserAvatarProps = {
  user?: User
  size: string
}

const UserAvatar = ({ user, size }: UserAvatarProps) => {
  const borderColor = useColorModeValue('gray.300', 'gray.700')
  const [borderColorDecoded] = useToken('colors', [borderColor])

  if (user) {
    return (
      <Avatar
        src={user.picture}
        width={size}
        height={size}
        border={`1px solid ${borderColorDecoded}`}
      />
    )
  } else {
    return <SkeletonCircle size={size} />
  }
}

export default UserAvatar

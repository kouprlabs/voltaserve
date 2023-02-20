import { Avatar, useColorModeValue, useToken } from '@chakra-ui/react'
import { User } from '@/api/user'

type AvatarImageProps = {
  user: User
}

const AvatarImage = ({ user }: AvatarImageProps) => {
  const borderColor = useColorModeValue('gray.300', 'gray.700')
  const [borderColorDecoded] = useToken('colors', [borderColor])
  return (
    <Avatar
      name={user.fullName}
      src={user.picture}
      size="sm"
      width="40px"
      height="40px"
      border={`1px solid ${borderColorDecoded}`}
    />
  )
}

export default AvatarImage

import { Avatar, useColorModeValue, useToken } from '@chakra-ui/react'
import { User } from '@/client/idp/user'

type AccountMenuAvatarImageProps = {
  user: User
}

const AccountMenuAvatarImage = ({ user }: AccountMenuAvatarImageProps) => {
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

export default AccountMenuAvatarImage

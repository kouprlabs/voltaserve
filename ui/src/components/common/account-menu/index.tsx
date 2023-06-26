import { Link } from 'react-router-dom'
import {
  HStack,
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Portal,
  SkeletonCircle,
  Stack,
  Text,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import ActiveCircle from './active-circle'
import AvatarButton from './avatar-button'
import AvatarImage from './avatar-image'

const AccountMenu = () => {
  const { data: user } = UserAPI.useGet(swrConfig())
  if (user) {
    return (
      <Menu>
        <MenuButton as={AvatarButton} user={user} />
        <Portal>
          <MenuList>
            <HStack spacing={variables.spacingXs} px={variables.spacingSm}>
              <AvatarImage user={user} />
              <Stack spacing={0}>
                <Text
                  fontWeight="semibold"
                  flexGrow={1}
                  textOverflow="ellipsis"
                  overflow="hidden"
                  whiteSpace="nowrap"
                >
                  {user.fullName}
                </Text>
                <Text color="gray.500">{user.email}</Text>
              </Stack>
            </HStack>
            <MenuDivider />
            <MenuItem as={Link} to="/account/settings">
              Account
            </MenuItem>
            <MenuItem as={Link} to="/sign-out" color="red">
              Sign Out
            </MenuItem>
          </MenuList>
        </Portal>
      </Menu>
    )
  } else {
    return (
      <ActiveCircle>
        <SkeletonCircle size="40px" />
      </ActiveCircle>
    )
  }
}

export default AccountMenu

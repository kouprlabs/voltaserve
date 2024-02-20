import { Link } from 'react-router-dom'
import {
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Portal,
  SkeletonCircle,
  Text,
} from '@chakra-ui/react'
import classNames from 'classnames'
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
            <div
              className={classNames(
                'flex',
                'flex-row',
                'items-center',
                'gap-0.5',
                'px-1',
              )}
            >
              <AvatarImage user={user} />
              <div className={classNames('flex', 'flex-col', 'gap-0')}>
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
              </div>
            </div>
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

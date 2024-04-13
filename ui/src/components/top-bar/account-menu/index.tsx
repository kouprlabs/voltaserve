import { Link } from 'react-router-dom'
import {
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Portal,
  SkeletonCircle,
} from '@chakra-ui/react'
import cx from 'classnames'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountMenuActiveCircle from './account-menu-active-circle'
import AccountMenuAvatarButton from './account-menu-avatar-button'
import AccountMenuAvatarImage from './account-menu-avatar-image'

const TopBarAccountMenu = () => {
  const { data: user } = UserAPI.useGet(swrConfig())
  if (user) {
    return (
      <Menu>
        <MenuButton as={AccountMenuAvatarButton} user={user} />
        <Portal>
          <MenuList>
            <div
              className={cx(
                'flex',
                'flex-row',
                'items-center',
                'gap-0.5',
                'px-1',
              )}
            >
              <AccountMenuAvatarImage user={user} />
              <div className={cx('flex', 'flex-col', 'gap-0')}>
                <span
                  className={cx(
                    'font-semibold',
                    'grow',
                    'text-ellipsis',
                    'overflow-hidden',
                    'whitespace-nowrap',
                  )}
                >
                  {user.fullName}
                </span>
                <span className={cx('text-gray-500')}>{user.email}</span>
              </div>
            </div>
            <MenuDivider />
            <MenuItem as={Link} to="/account/settings">
              Account
            </MenuItem>
            <MenuItem as={Link} to="/sign-out" className={cx('text-red-500')}>
              Sign Out
            </MenuItem>
          </MenuList>
        </Portal>
      </Menu>
    )
  } else {
    return (
      <AccountMenuActiveCircle>
        <SkeletonCircle size="40px" />
      </AccountMenuActiveCircle>
    )
  }
}

export default TopBarAccountMenu

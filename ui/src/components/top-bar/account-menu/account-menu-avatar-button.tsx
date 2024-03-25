import { Avatar, useColorModeValue, useToken } from '@chakra-ui/react'
import { forwardRef } from '@chakra-ui/system'
import classNames from 'classnames'
import { User } from '@/client/idp/user'
import { useAppSelector } from '@/store/hook'
import { NavType } from '@/store/ui/nav'
import AccountMenuActiveCircle from './account-menu-active-circle'

type AccountMenuAvatarButtonProps = {
  user: User
}

const AccountMenuAvatarButton = forwardRef<AccountMenuAvatarButtonProps, 'div'>(
  ({ user, ...props }, ref) => {
    const borderColor = useToken(
      'colors',
      useColorModeValue('gray.300', 'gray.700'),
    )
    const activeNav = useAppSelector((state) => state.ui.nav.active)
    return (
      <div ref={ref} {...props} className={classNames('cursor-pointer')}>
        <AccountMenuActiveCircle>
          <Avatar
            name={user.fullName}
            src={user.picture}
            size="sm"
            width="40px"
            height="40px"
            border={
              activeNav === NavType.Account
                ? 'none'
                : `1px solid ${borderColor}`
            }
          />
        </AccountMenuActiveCircle>
      </div>
    )
  },
)

export default AccountMenuAvatarButton

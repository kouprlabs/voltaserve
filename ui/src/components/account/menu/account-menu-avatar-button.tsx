import { Avatar } from '@chakra-ui/react'
import { forwardRef } from '@chakra-ui/system'
import cx from 'classnames'
import InvitationAPI from '@/client/api/invitation'
import { User } from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import NotificationBadge from '@/lib/components/notification-badge'
import { useAppSelector } from '@/store/hook'
import { NavType } from '@/store/ui/nav'
import AccountMenuActiveCircle from './account-menu-active-circle'

export type AccountMenuAvatarButtonProps = {
  user: User
}

const AccountMenuAvatarButton = forwardRef<AccountMenuAvatarButtonProps, 'div'>(
  ({ user, ...props }, ref) => {
    const activeNav = useAppSelector((state) => state.ui.nav.active)
    const { data: count } = InvitationAPI.useGetIncomingCount(swrConfig())
    const isActive = activeNav === NavType.Account

    return (
      <NotificationBadge hasBadge={count && count > 0 ? true : false}>
        <div ref={ref} {...props} className={cx('cursor-pointer')}>
          <AccountMenuActiveCircle>
            <Avatar
              name={user.fullName}
              src={user.picture}
              size="sm"
              className={cx('w-[40px]', 'h-[40px]', {
                'border': isActive,
                'border-gray-300': isActive,
                'dark:border-gray-700': isActive,
              })}
            />
          </AccountMenuActiveCircle>
        </div>
      </NotificationBadge>
    )
  },
)

export default AccountMenuAvatarButton

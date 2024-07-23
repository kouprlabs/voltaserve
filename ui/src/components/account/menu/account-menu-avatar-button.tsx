// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
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
      <div ref={ref} {...props} className={cx('cursor-pointer')}>
        <AccountMenuActiveCircle>
          <NotificationBadge hasBadge={count && count > 0 ? true : false}>
            <Avatar
              name={user.fullName}
              src={user.picture}
              size="sm"
              className={cx(
                'w-[40px]',
                'h-[40px]',
                'border',
                'border-gray-300',
                {
                  'dark:border-gray-700': !isActive,
                },
              )}
            />
          </NotificationBadge>
        </AccountMenuActiveCircle>
      </div>
    )
  },
)

export default AccountMenuAvatarButton

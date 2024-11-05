// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Link } from 'react-router-dom'
import { MenuItem } from '@chakra-ui/react'
import { AccountMenu as KouprAccountMenu, NumberTag } from '@koupr/ui'
import cx from 'classnames'
import InvitationAPI from '@/client/api/invitation'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import { getPictureUrl } from '@/lib/helpers/picture'

const AccountMenu = () => {
  const { data: user } = UserAPI.useGet(swrConfig())
  const { data: invitationCount } =
    InvitationAPI.useGetIncomingCount(swrConfig())

  return (
    <KouprAccountMenu
      isLoading={!user}
      hasBadge={Boolean(invitationCount && invitationCount > 0)}
      name={user?.fullName}
      email={user?.email}
      picture={user?.picture ? getPictureUrl(user.picture) : undefined}
      menuItems={
        <>
          <MenuItem as={Link} to="/account/settings">
            Settings
          </MenuItem>
          <MenuItem as={Link} to="/account/invitation">
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
              <span>Invitations</span>
              {invitationCount && invitationCount > 0 ? (
                <NumberTag>{invitationCount}</NumberTag>
              ) : null}
            </div>
          </MenuItem>
          <MenuItem as={Link} to="/sign-out" className={cx('text-red-500')}>
            Sign Out
          </MenuItem>
        </>
      }
    />
  )
}

export default AccountMenu

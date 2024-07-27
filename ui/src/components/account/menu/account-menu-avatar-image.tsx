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
import cx from 'classnames'
import { User } from '@/client/idp/user'

export type AccountMenuAvatarImageProps = {
  user: User
}

const AccountMenuAvatarImage = ({ user }: AccountMenuAvatarImageProps) => (
  <Avatar
    name={user.fullName}
    src={user.picture}
    size="sm"
    className={cx(
      'w-[40px]',
      'h-[40px]',
      'border',
      'border-gray-300',
      'dark:border-gray-700',
    )}
  />
)

export default AccountMenuAvatarImage

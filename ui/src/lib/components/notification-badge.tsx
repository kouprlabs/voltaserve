// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ReactNode } from 'react'
import { Circle } from '@chakra-ui/react'
import cx from 'classnames'

export type NotificationBadgeProps = {
  hasBadge?: boolean
  children?: ReactNode
}

const NotificationBadge = ({ hasBadge, children }: NotificationBadgeProps) => {
  return (
    <div className={cx('flex', 'items-center', 'justify-center', 'relative')}>
      {children}
      {hasBadge ? (
        <Circle size="10px" bg="red" position="absolute" top={0} right={0} />
      ) : null}
    </div>
  )
}

export default NotificationBadge

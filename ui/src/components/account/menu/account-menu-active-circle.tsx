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
import variables from '@/lib/variables'
import { useAppSelector } from '@/store/hook'
import { NavType } from '@/store/ui/nav'

export type AccountMenuActiveCircleProps = {
  children?: ReactNode
}

const AccountMenuActiveCircle = ({
  children,
}: AccountMenuActiveCircleProps) => {
  const activeNav = useAppSelector((state) => state.ui.nav.active)
  return (
    <Circle
      className={cx('w-[50px]', 'h-[50px]')}
      bg={activeNav === NavType.Account ? variables.gradiant : 'transparent'}
    >
      {children}
    </Circle>
  )
}

export default AccountMenuActiveCircle

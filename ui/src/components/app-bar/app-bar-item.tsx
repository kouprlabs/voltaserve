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
import { Link as ChakraLink } from '@chakra-ui/react'
import cx from 'classnames'

export type AppBarItemProps = {
  title: string
  href: string
  isActive: boolean
}

const AppBarItem = ({ title, href, isActive }: AppBarItemProps) => (
  <ChakraLink
    as={Link}
    to={href}
    lineHeight="40px"
    variant="no-underline"
    className={cx(
      'opacity-100',
      'hover:opacity-80',
      'h-[40px]',
      'rounded-xl',
      'pt-0',
      'pr-[20px]',
      'pb-0',
      'pl-[20px]',
      'font-semibold',
      {
        'text-white': isActive,
        'dark:text-gray-600': isActive,
        'bg-black': isActive,
        'dark:bg-white': isActive,
        'bg-transparent': !isActive,
      },
    )}
  >
    {title}
  </ChakraLink>
)

export default AppBarItem

// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { ReactNode, useMemo } from 'react'
import { Tag, useColorMode } from '@chakra-ui/react'
import cx from 'classnames'

export type TabTagProps = {
  children?: ReactNode
  className?: string
  isActive?: boolean
}

const NumberTag = ({ isActive, className, children }: TabTagProps) => {
  const { colorMode } = useColorMode()
  const bg = useMemo(() => {
    if (isActive) {
      if (colorMode === 'light') {
        return 'white'
      } else if (colorMode === 'dark') {
        return 'gray.800'
      }
    } else {
      if (colorMode === 'light') {
        return 'black'
      } else if (colorMode === 'dark') {
        return 'white'
      }
    }
  }, [isActive, colorMode])
  const color = useMemo(() => {
    if (isActive) {
      if (colorMode === 'light') {
        return 'black'
      } else if (colorMode === 'dark') {
        return 'white'
      }
    } else {
      if (colorMode === 'light') {
        return 'white'
      } else if (colorMode === 'dark') {
        return 'gray.800'
      }
    }
  }, [isActive, colorMode])

  return (
    <Tag className={cx('rounded-full', className)} color={color} bg={bg}>
      {children}
    </Tag>
  )
}

export default NumberTag

// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Circle, Tooltip } from '@chakra-ui/react'
import { IconVisibility } from '@koupr/ui'
import cx from 'classnames'

const IconBadgeInsights = () => (
  <Tooltip label="This item has insights">
    <Circle
      className={cx(
        'text-orange-600',
        'bg-white',
        'w-[23px]',
        'h-[23px]',
        'border',
        'border-gray-200',
      )}
    >
      <IconVisibility className={cx('text-[12px]')} />
    </Circle>
  </Tooltip>
)

export default IconBadgeInsights

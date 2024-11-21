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
import { IconModeHeat } from '@koupr/ui'
import cx from 'classnames'

const IconBadgeMosaic = () => (
  <Tooltip label="This item has a mosaic">
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
      <IconModeHeat className={cx('text-[12px]')} />
    </Circle>
  </Tooltip>
)

export default IconBadgeMosaic

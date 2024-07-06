// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useCallback } from 'react'
import { CenterProps, useColorMode } from '@chakra-ui/react'
import cx from 'classnames'
import LogoDarkGlossySvg from './logo-dark-glossy.svg?react'
import LogoDarkSvg from './logo-dark.svg?react'
import LogoGlossySvg from './logo-glossy.svg?react'
import LogoSvg from './logo.svg?react'

type LogoProps = CenterProps & {
  isGlossy?: boolean
}

const Logo = ({ isGlossy = false }: LogoProps) => {
  const { colorMode } = useColorMode()
  const renderSvg = useCallback(() => {
    if (isGlossy) {
      return colorMode === 'dark' ? <LogoDarkGlossySvg /> : <LogoGlossySvg />
    } else {
      return colorMode === 'dark' ? <LogoDarkSvg /> : <LogoSvg />
    }
  }, [colorMode, isGlossy])
  return (
    <div className={cx('flex', 'items-center', 'justify-center')}>
      {renderSvg()}
    </div>
  )
}

export default Logo

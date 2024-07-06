// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import cx from 'classnames'
import { FileCommonProps } from '@/types/file'
import IconBadge from '../icon-badge'
import FolderSvg from './assets/icon-folder.svg'

type IconFolderProps = {
  isLoading?: boolean
} & FileCommonProps

const MIN_WIDTH = 45
const MIN_HEIGHT = 36.05
const BASE_WIDTH = 67
const BASE_HEIGHT = 53.68

const IconFolder = ({ file, scale, isLoading }: IconFolderProps) => {
  const width = Math.max(MIN_WIDTH, BASE_WIDTH * scale)
  const height = Math.max(MIN_HEIGHT, BASE_HEIGHT * scale)

  return (
    <>
      <img
        src={FolderSvg}
        className="pointer-events-none select-none"
        style={{ width: `${width}px`, height: `${height}px` }}
      />
      <div
        className={cx(
          'absolute',
          'flex',
          'flex-row',
          'items-center',
          'gap-[2px]',
          'bottom-[-5px]',
          'right-[-5px]',
        )}
      >
        <IconBadge file={file} isLoading={isLoading} />
      </div>
    </>
  )
}

export default IconFolder

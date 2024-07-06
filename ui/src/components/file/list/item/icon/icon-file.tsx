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
import IconBadge from './icon-badge'
import IconDiverse from './icon-diverse'
import IconThumbnail from './icon-thumbnail'

type IconFileProps = FileCommonProps

const IconFile = ({ file, scale }: IconFileProps) => (
  <>
    {file.snapshot?.thumbnail ? (
      <IconThumbnail file={file} scale={scale} />
    ) : (
      <>
        <IconDiverse file={file} scale={scale} />
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
          <IconBadge file={file} />
        </div>
      </>
    )}
  </>
)

export default IconFile

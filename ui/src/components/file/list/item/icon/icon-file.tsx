// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import cx from 'classnames'
import { FileCommonProps } from '@/types/file'
import IconBadge from './icon-badge'
import IconDiverse from './icon-diverse'
import IconThumbnail from './icon-thumbnail'

type IconFileProps = FileCommonProps

const IconFile = ({ file, scale, isLoading }: IconFileProps) => (
  <>
    {file.snapshot?.thumbnail ? (
      <IconThumbnail file={file} scale={scale} isLoading={isLoading} />
    ) : (
      <>
        <IconDiverse file={file} scale={scale} />
        <div
          className={cx('absolute', 'flex', 'flex-row', 'items-center', 'gap-[2px]', 'bottom-[-5px]', 'right-[-5px]')}
        >
          <IconBadge file={file} isLoading={isLoading} />
        </div>
      </>
    )}
  </>
)

export default IconFile

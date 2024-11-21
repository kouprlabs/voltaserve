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
import { computeScale } from '../scale'
import IconFile from './icon-file'
import IconFolder from './icon-folder'

export type ItemIconProps = {
  isLoading?: boolean
} & FileCommonProps

const ItemIcon = ({ file, scale, viewType, isLoading }: ItemIconProps) => (
  <>
    <div
      className={cx('z-0', 'text-gray-500', 'dark:text-gray-300', 'relative')}
    >
      {file.type === 'file' ? (
        <IconFile
          file={file}
          scale={computeScale(scale, viewType)}
          viewType={viewType}
          isLoading={isLoading}
        />
      ) : file.type === 'folder' ? (
        <IconFolder
          file={file}
          scale={computeScale(scale, viewType)}
          viewType={viewType}
          isLoading={isLoading}
        />
      ) : null}
    </div>
  </>
)

export default ItemIcon

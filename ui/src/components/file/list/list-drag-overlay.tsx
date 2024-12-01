// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { DragOverlay } from '@dnd-kit/core'
import cx from 'classnames'
import { useAppSelector } from '@/store/hook'
import { FileCommonProps } from '@/types/file'
import ListItem from './item'

type ListDragOverlayProps = FileCommonProps

const ListDragOverlay = ({ file, scale, viewType }: ListDragOverlayProps) => {
  const selectionCount = useAppSelector((state) => state.ui.files.selection.length)

  return (
    <DragOverlay>
      <div className={cx('relative')}>
        <ListItem file={file} scale={scale} isPresentational={true} isDragging={true} viewType={viewType} />
        {selectionCount > 1 ? (
          <div
            className={cx(
              'absolute',
              'flex',
              'items-center',
              'justify-center',
              'bottom-[-5px]',
              'right-[-5px]',
              'text-white',
              'rounded-xl',
              'min-w-[30px]',
              'h-[30px]',
              'px-1',
              'bg-blue-500',
            )}
          >
            {selectionCount}
          </div>
        ) : null}
      </div>
    </DragOverlay>
  )
}

export default ListDragOverlay

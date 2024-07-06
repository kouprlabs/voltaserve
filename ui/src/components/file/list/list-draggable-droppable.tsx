// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useState, MouseEvent } from 'react'
import {
  DragCancelEvent,
  DragEndEvent,
  DragStartEvent,
  useDndMonitor,
  useDraggable,
  useDroppable,
} from '@dnd-kit/core'
import cx from 'classnames'
import FileAPI, { FileType } from '@/client/api/file'
import store from '@/store/configure-store'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { hiddenUpdated, selectionUpdated } from '@/store/ui/files'
import { FileCommonProps } from '@/types/file'
import ListItem from './item'

export type ListDraggableDroppableProps = {
  onContextMenu?: (event: MouseEvent) => void
} & FileCommonProps

const ListDraggableDroppable = ({
  file,
  scale,
  viewType,
  isSelectionMode,
  onContextMenu,
}: ListDraggableDroppableProps) => {
  const dispatch = useAppDispatch()
  const [isVisible, setVisible] = useState(true)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const {
    attributes,
    listeners,
    setNodeRef: setDraggableNodeRef,
  } = useDraggable({
    id: file.id,
  })
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const { isOver, setNodeRef: setDroppableNodeRef } = useDroppable({
    id: file.id,
  })
  const [isLoading, setIsLoading] = useState(false)

  useDndMonitor({
    onDragStart: (event: DragStartEvent) => {
      if (selection.includes(file.id) || event.active.id === file.id) {
        setVisible(false)
      }
    },
    onDragEnd: async (event: DragEndEvent) => {
      if (selection.includes(file.id) || event.active.id === file.id) {
        setVisible(true)
      }
      if (
        file.type === FileType.Folder &&
        file.id !== event.active.id &&
        !selection.includes(file.id) &&
        isOver
      ) {
        const idsToMove = [
          ...new Set<string>([...selection, event.active.id as string]),
        ]
        const list = store.getState().entities.files.list
        if (list) {
          mutateList?.({
            ...list,
            data: list.data.filter((e) => !idsToMove.includes(e.id)),
          })
        }
        dispatch(hiddenUpdated(idsToMove))
        setIsLoading(true)
        try {
          await FileAPI.move(file.id, { ids: idsToMove })
          mutateList?.()
        } finally {
          setIsLoading(false)
          dispatch(hiddenUpdated([]))
          dispatch(selectionUpdated([]))
        }
      }
    },
    onDragCancel: (event: DragCancelEvent) => {
      if (selection.includes(file.id) || event.active.id === file.id) {
        setVisible(true)
      }
    },
  })

  return (
    <>
      {file.type === FileType.File ? (
        <div
          ref={setDraggableNodeRef}
          className={cx(
            'border-2',
            'border-transparent',
            'hover:outline-none',
            'focus:outline-none',
            'cursor-default',
            { 'visible': isVisible },
            { 'invisible': !isVisible },
          )}
          {...listeners}
          {...attributes}
        >
          <ListItem
            file={file}
            scale={scale}
            viewType={viewType}
            isSelectionMode={isSelectionMode}
            onContextMenu={onContextMenu}
          />
        </div>
      ) : null}
      {file.type === FileType.Folder ? (
        <div
          ref={setDraggableNodeRef}
          className={cx(
            'border-2',
            'rounded-md',
            'hover:outline-none',
            'focus:outline-none',
            { 'visible': isVisible },
            { 'invisible': !isVisible },
            { 'border-blue-500': isOver },
            {
              'bg-blue-100': isOver,
              'dark:bg-blue-950': isOver,
            },
            { 'border-transparent': !isOver },
          )}
          {...listeners}
          {...attributes}
        >
          <div ref={setDroppableNodeRef}>
            <ListItem
              file={file}
              scale={scale}
              viewType={viewType}
              isPresentational={isOver}
              isLoading={isLoading}
              isSelectionMode={isSelectionMode}
              onContextMenu={onContextMenu}
            />
          </div>
        </div>
      ) : null}
    </>
  )
}

export default ListDraggableDroppable

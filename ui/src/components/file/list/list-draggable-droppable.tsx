import { useState, MouseEvent } from 'react'
import { useParams } from 'react-router-dom'
import { Box, useToken } from '@chakra-ui/react'
import { useSWRConfig } from 'swr'
import {
  DragCancelEvent,
  DragEndEvent,
  DragStartEvent,
  useDndMonitor,
  useDraggable,
  useDroppable,
} from '@dnd-kit/core'
import classNames from 'classnames'
import FileAPI, { FileType, List } from '@/client/api/file'
import useFileListSearchParams from '@/hooks/use-file-list-params'
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
  const { mutate } = useSWRConfig()
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const [isVisible, setVisible] = useState(true)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const {
    attributes,
    listeners,
    setNodeRef: setDraggableNodeRef,
  } = useDraggable({
    id: file.id,
  })
  const { isOver, setNodeRef: setDroppableNodeRef } = useDroppable({
    id: file.id,
  })
  const [isLoading, setIsLoading] = useState(false)
  const fileListSearchParams = useFileListSearchParams()
  const borderColor = useToken('colors', 'green.300')

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
          await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`, {
            ...list,
            data: list.data.filter((e) => !idsToMove.includes(e.id)),
          })
        }
        dispatch(hiddenUpdated(idsToMove))
        setIsLoading(true)
        await FileAPI.move(file.id, { ids: idsToMove })
        await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
        setIsLoading(false)
        dispatch(hiddenUpdated([]))
        dispatch(selectionUpdated([]))
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
        <Box
          ref={setDraggableNodeRef}
          className={classNames(
            'border-2',
            'border-transparent',
            { 'visible': isVisible },
            { 'hidden': !isVisible },
          )}
          _hover={{ outline: 'none' }}
          _focus={{ outline: 'none' }}
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
        </Box>
      ) : null}
      {file.type === FileType.Folder ? (
        <Box
          ref={setDraggableNodeRef}
          className={classNames(
            'border-2',
            'rounded-md',
            { 'visible': isVisible },
            { 'invisible': !isVisible },
          )}
          style={{ borderColor: isOver ? borderColor : 'transparent' }}
          _hover={{ outline: 'none' }}
          _focus={{ outline: 'none' }}
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
        </Box>
      ) : null}
    </>
  )
}

export default ListDraggableDroppable

import { useState, MouseEvent } from 'react'
import { useParams } from 'react-router-dom'
import { Box } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { useSWRConfig } from 'swr'
import {
  DragCancelEvent,
  DragEndEvent,
  DragStartEvent,
  useDndMonitor,
  useDraggable,
  useDroppable,
} from '@dnd-kit/core'
import FileAPI, { File, FileType } from '@/client/api/file'
import store from '@/store/configure-store'
import { filesRemoved } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { hiddenItemsUpdated, selectedItemsUpdated } from '@/store/ui/files'
import Item from './item'

type ItemDraggableDroppableProps = {
  file: File
  scale: number
  onContextMenu?: (event: MouseEvent) => void
}

const ItemDraggableDroppable = ({
  file,
  scale,
  onContextMenu,
}: ItemDraggableDroppableProps) => {
  const { mutate } = useSWRConfig()
  const dispatch = useAppDispatch()
  const params = useParams()
  const fileId = params.fileId as string
  const [isVisible, setVisible] = useState(true)
  const selectedItems = useAppSelector((state) => state.ui.files.selectedItems)
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

  useDndMonitor({
    onDragStart: (event: DragStartEvent) => {
      if (selectedItems.includes(file.id) || event.active.id === file.id) {
        setVisible(false)
      }
    },
    onDragEnd: async (event: DragEndEvent) => {
      if (selectedItems.includes(file.id) || event.active.id === file.id) {
        setVisible(true)
      }
      if (
        file.type === FileType.Folder &&
        file.id !== event.active.id &&
        !selectedItems.includes(file.id) &&
        isOver
      ) {
        const idsToMove = [
          ...new Set<string>([...selectedItems, event.active.id as string]),
        ]
        const list = store.getState().entities.files.list
        if (list) {
          await mutate(`/files/${fileId}/list`, {
            ...list,
            data: list.data.filter((e) => !idsToMove.includes(e.id)),
          })
        }
        dispatch(filesRemoved({ id: fileId, files: idsToMove }))
        dispatch(hiddenItemsUpdated(idsToMove))
        setIsLoading(true)
        await FileAPI.move(file.id, { ids: idsToMove })
        await mutate(`/files/${fileId}/list`)
        setIsLoading(false)
        dispatch(hiddenItemsUpdated([]))
        dispatch(selectedItemsUpdated([]))
      }
    },
    onDragCancel: (event: DragCancelEvent) => {
      if (selectedItems.includes(file.id) || event.active.id === file.id) {
        setVisible(true)
      }
    },
  })

  return (
    <>
      {file.type === FileType.File ? (
        <Box
          ref={setDraggableNodeRef}
          border="2px solid"
          borderColor="transparent"
          visibility={isVisible ? 'visible' : 'hidden'}
          _hover={{ outline: 'none' }}
          _focus={{ outline: 'none' }}
          {...listeners}
          {...attributes}
        >
          <Item file={file} scale={scale} onContextMenu={onContextMenu} />
        </Box>
      ) : null}
      {file.type === FileType.Folder ? (
        <Box
          ref={setDraggableNodeRef}
          border="2px solid"
          borderColor={isOver ? 'green.300' : 'transparent'}
          borderRadius={variables.borderRadiusSm}
          visibility={isVisible ? 'visible' : 'hidden'}
          _hover={{ outline: 'none' }}
          _focus={{ outline: 'none' }}
          {...listeners}
          {...attributes}
        >
          <Box ref={setDroppableNodeRef}>
            <Item
              file={file}
              scale={scale}
              isPresentational={isOver}
              isLoading={isLoading}
              onContextMenu={onContextMenu}
            />
          </Box>
        </Box>
      ) : null}
    </>
  )
}

export default ItemDraggableDroppable

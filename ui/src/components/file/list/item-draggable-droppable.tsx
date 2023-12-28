import { ReactNode, useState } from 'react'
import { Box } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import {
  DragCancelEvent,
  DragEndEvent,
  DragStartEvent,
  useDndMonitor,
  useDraggable,
  useDroppable,
} from '@dnd-kit/core'
import { File, FileType } from '@/client/api/file'
import { useAppSelector } from '@/store/hook'

type ItemDraggableDroppableProps = {
  file: File
  children?: ReactNode
}

const ItemDraggableDroppable = ({
  file,
  children,
}: ItemDraggableDroppableProps) => {
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
  useDndMonitor({
    onDragStart: (event: DragStartEvent) => {
      if (selection.includes(file.id) || event.active.id === file.id) {
        setVisible(false)
      }
    },
    onDragEnd: (event: DragEndEvent) => {
      if (selection.includes(file.id) || event.active.id === file.id) {
        setVisible(true)
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
          visibility={isVisible ? 'visible' : 'hidden'}
          {...listeners}
          {...attributes}
        >
          {children}
        </Box>
      ) : null}
      {file.type === FileType.Folder ? (
        <Box
          ref={setDraggableNodeRef}
          border="2px solid"
          borderColor={isOver ? 'purple.200' : 'transparent'}
          borderRadius={variables.borderRadiusSm}
          visibility={isVisible ? 'visible' : 'hidden'}
          {...listeners}
          {...attributes}
        >
          <Box ref={setDroppableNodeRef}>{children}</Box>
        </Box>
      ) : null}
    </>
  )
}

export default ItemDraggableDroppable

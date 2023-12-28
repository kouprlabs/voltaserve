import { ReactNode } from 'react'
import { Box } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { useDraggable, useDroppable } from '@dnd-kit/core'
import { File, FileType } from '@/client/api/file'

type DraggableDroppableProps = {
  file: File
  children?: ReactNode
}

const DraggableDroppable = ({ file, children }: DraggableDroppableProps) => {
  const {
    attributes,
    listeners,
    setNodeRef: setDraggableNodeRef,
    transform,
    isDragging,
  } = useDraggable({
    id: file.id,
  })
  const { isOver, setNodeRef: setDroppableNodeRef } = useDroppable({
    id: file.id,
  })
  const draggableStyle = transform
    ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
      }
    : undefined
  const droppableBorderColor = () => {
    if (isOver && !isDragging) {
      return '#58D68D'
    } else if (isDragging) {
      return '#DE3163'
    } else {
      return 'transparent'
    }
  }

  return (
    <>
      {file.type === FileType.File ? (
        <Box
          ref={setDraggableNodeRef}
          style={draggableStyle}
          border="2px solid"
          borderColor={isDragging ? '#DE3163' : 'transparent'}
          borderRadius={variables.borderRadiusSm}
          zIndex={isDragging ? 999999 : 'auto'}
          {...listeners}
          {...attributes}
        >
          {children}
        </Box>
      ) : null}
      {file.type === FileType.Folder ? (
        <Box
          ref={setDraggableNodeRef}
          style={draggableStyle}
          border="2px solid"
          borderColor={droppableBorderColor()}
          borderRadius={variables.borderRadiusSm}
          zIndex={isDragging ? 999999 : 'auto'}
          {...listeners}
          {...attributes}
        >
          <Box ref={setDroppableNodeRef}>{children}</Box>
        </Box>
      ) : null}
    </>
  )
}

export default DraggableDroppable

import { Box, Badge } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { DragOverlay } from '@dnd-kit/core'
import { File } from '@/client/api/file'
import { useAppSelector } from '@/store/hook'
import Item from './item'

type ItemDragOverlayProps = {
  file: File
  scale: number
}

const ItemDragOverlay = ({ file, scale }: ItemDragOverlayProps) => {
  const selectionCount = useAppSelector(
    (state) => state.ui.files.selection.length,
  )

  return (
    <DragOverlay>
      <Box position="relative">
        <Item file={file} scale={scale} />
        {selectionCount > 1 ? (
          <Badge
            position="absolute"
            bottom="-5px"
            right="-5px"
            colorScheme="green"
            borderRadius="30px"
            minW="30px"
            h="30px"
            px={variables.spacingSm}
            display="flex"
            alignItems="center"
            justifyContent="center"
          >
            {selectionCount}
          </Badge>
        ) : null}
      </Box>
    </DragOverlay>
  )
}

export default ItemDragOverlay

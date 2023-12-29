import { Box, Center } from '@chakra-ui/react'
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
        <Item file={file} scale={scale} isPresentational={true} />
        {selectionCount > 1 ? (
          <Center
            position="absolute"
            bottom="-5px"
            right="-5px"
            color="white"
            bgColor="green.300"
            borderRadius="30px"
            minW="30px"
            h="30px"
            px={variables.spacingSm}
            display="flex"
            alignItems="center"
            justifyContent="center"
          >
            {selectionCount}
          </Center>
        ) : null}
      </Box>
    </DragOverlay>
  )
}

export default ItemDragOverlay

import { useMemo } from 'react'
import { Box, HStack } from '@chakra-ui/react'
import { FcFolder } from 'react-icons/fc'
import { File } from '@/client/api/file'
import ProcessingBadge from './processing-badge'
import SharedBadge from './shared-badge'

type FolderIconProps = {
  file: File
  scale: number
  isLoading?: boolean
}

const ICON_FONT_SIZE = 92

const FolderIcon = ({ file, scale, isLoading }: FolderIconProps) => {
  const fontSize = useMemo(() => `${ICON_FONT_SIZE * scale}px`, [scale])
  return (
    <Box position="relative">
      <FcFolder fontSize={fontSize} />
      <HStack position="absolute" bottom="7px" right="2px" spacing="2px">
        {file.isShared && <SharedBadge />}
        {isLoading && <ProcessingBadge />}
      </HStack>
    </Box>
  )
}

export default FolderIcon

import { useMemo } from 'react'
import { Box, HStack } from '@chakra-ui/react'
import { FcFolder } from 'react-icons/fc'
import { File } from '@/client/api/file'
import SharedSign from './shared-sign'

type FolderIconProps = {
  file: File
  scale: number
}

const ICON_FONT_SIZE = 92

const FolderIcon = ({ file, scale }: FolderIconProps) => {
  const fontSize = useMemo(() => `${ICON_FONT_SIZE * scale}px`, [scale])
  return (
    <Box position="relative">
      <FcFolder fontSize={fontSize} />
      <HStack position="absolute" bottom="7px" right="2px" spacing="2px">
        {file.isShared && <SharedSign />}
      </HStack>
    </Box>
  )
}

export default FolderIcon

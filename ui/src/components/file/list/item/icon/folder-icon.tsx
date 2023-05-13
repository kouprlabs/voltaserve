import { useMemo } from 'react'
import { Box } from '@chakra-ui/react'
import { FcFolder } from 'react-icons/fc'
import { File } from '@/api/file'
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
      {file.isShared && <SharedSign bottom="7px" right="2px" />}
    </Box>
  )
}

export default FolderIcon

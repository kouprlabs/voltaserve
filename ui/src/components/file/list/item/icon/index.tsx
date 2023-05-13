import { Box, useColorModeValue } from '@chakra-ui/react'
import { File } from '@/api/file'
import { ItemSize } from '..'
import FileIcon from './file-icon'
import FolderIcon from './folder-icon'

type IconProps = {
  file: File
  size: ItemSize
}

const Icon = ({ file, size }: IconProps) => {
  const color = useColorModeValue('gray.500', 'gray.300')
  if (file.type === 'file') {
    return (
      <Box color={color} zIndex={0}>
        <FileIcon file={file} size={size} />
      </Box>
    )
  } else if (file.type === 'folder') {
    return (
      <Box zIndex={0}>
        <FolderIcon file={file} size={size} />
      </Box>
    )
  } else {
    return null
  }
}

export default Icon

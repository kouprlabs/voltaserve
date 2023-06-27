import { Box, useColorModeValue } from '@chakra-ui/react'
import { File } from '@/client/api/file'
import FileIcon from './file-icon'
import FolderIcon from './folder-icon'

type IconProps = {
  file: File
  scale: number
}

const Icon = ({ file, scale }: IconProps) => {
  const color = useColorModeValue('gray.500', 'gray.300')
  if (file.type === 'file') {
    return (
      <Box color={color} zIndex={0}>
        <FileIcon file={file} scale={scale} />
      </Box>
    )
  } else if (file.type === 'folder') {
    return (
      <Box zIndex={0}>
        <FolderIcon file={file} scale={scale} />
      </Box>
    )
  } else {
    return null
  }
}

export default Icon

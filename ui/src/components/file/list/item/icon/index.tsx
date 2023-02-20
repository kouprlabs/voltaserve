import { Box, useColorModeValue } from '@chakra-ui/react'
import { File } from '@/api/file'
import { ItemSize } from '..'
import FileListItemFileIcon from './file-icon'
import FileListItemFolderIcon from './folder-icon'

type FileListItemIconProps = {
  file: File
  size: ItemSize
}

const FileListItemIcon = ({ file, size }: FileListItemIconProps) => {
  const color = useColorModeValue('gray.500', 'gray.300')
  if (file.type === 'file') {
    return (
      <Box color={color} zIndex={0}>
        <FileListItemFileIcon file={file} size={size} />
      </Box>
    )
  } else if (file.type === 'folder') {
    return (
      <Box zIndex={0}>
        <FileListItemFolderIcon file={file} size={size} />
      </Box>
    )
  } else {
    return null
  }
}

export default FileListItemIcon

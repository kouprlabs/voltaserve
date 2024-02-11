import { Box, useColorModeValue } from '@chakra-ui/react'
import { CommonItemProps } from '@/types/file'
import FileIcon from './file-icon'
import FolderIcon from './folder-icon'

type IconProps = {
  isLoading?: boolean
} & CommonItemProps

const Icon = ({ file, scale, viewType, isLoading }: IconProps) => {
  const color = useColorModeValue('gray.500', 'gray.300')
  if (file.type === 'file') {
    return (
      <Box color={color} zIndex={0}>
        <FileIcon file={file} scale={scale} viewType={viewType} />
      </Box>
    )
  } else if (file.type === 'folder') {
    return (
      <Box color={color} zIndex={0}>
        <FolderIcon
          file={file}
          scale={scale}
          viewType={viewType}
          isLoading={isLoading}
        />
      </Box>
    )
  } else {
    return null
  }
}

export default Icon

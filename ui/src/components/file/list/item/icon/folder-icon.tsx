import { useMemo } from 'react'
import { Box } from '@chakra-ui/react'
import { FcFolder } from 'react-icons/fc'
import { File } from '@/api/file'
import { ItemSize } from '..'
import SharedSign from './shared-sign'

type FolderIconProps = {
  file: File
  size: ItemSize
}

const FolderIcon = ({ file, size }: FolderIconProps) => {
  const fontSize = useMemo(() => {
    if (size === 'normal') {
      return '92px'
    }
    if (size === 'large') {
      return '150px'
    }
  }, [size])

  return (
    <Box position="relative">
      <FcFolder fontSize={fontSize} />
      {file.isShared && <SharedSign bottom="7px" right="2px" />}
    </Box>
  )
}

export default FolderIcon

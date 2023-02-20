import { useMemo } from 'react'
import { Box } from '@chakra-ui/react'
import { FcFolder } from 'react-icons/fc'
import { File } from '@/api/file'
import { ItemSize } from '..'
import FileListItemSharedSign from './shared-sign'

type FileListItemFolderIconProps = {
  file: File
  size: ItemSize
}

const FileListItemFolderIcon = ({
  file,
  size,
}: FileListItemFolderIconProps) => {
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
      {file.isShared && <FileListItemSharedSign bottom="7px" right="2px" />}
    </Box>
  )
}

export default FileListItemFolderIcon

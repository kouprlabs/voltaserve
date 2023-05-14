import { useMemo } from 'react'
import { Box } from '@chakra-ui/react'
import { File } from '@/api/file'
import * as fileExtension from '@/helpers/file-extension'
import FontIcon from './font-icon'
import ImageIcon from './image-icon'
import SharedSign from './shared-sign'
import Thumbnail from './thumbnail'

type FileIconProps = {
  file: File
  scale: number
}

const FileIcon = ({ file, scale }: FileIconProps) => {
  const isImage = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isImage(file.original.extension),
    [file.original]
  )
  if (isImage) {
    return <ImageIcon file={file} scale={scale} />
  } else {
    if (file.thumbnail) {
      return <Thumbnail file={file} scale={scale} />
    } else {
      return (
        <Box position="relative">
          <FontIcon file={file} scale={scale} />
          {file.isShared && <SharedSign bottom="-5px" right="0px" />}
        </Box>
      )
    }
  }
}

export default FileIcon

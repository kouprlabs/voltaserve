import { Box, HStack } from '@chakra-ui/react'
import { File } from '@/client/api/file'
import FontIcon from './font-icon'
import OcrBadge from './ocr-badge'
import SharedBadge from './shared-badge'
import Thumbnail from './thumbnail'

type FileIconProps = {
  file: File
  scale: number
}

const FileIcon = ({ file, scale }: FileIconProps) => {
  if (file.thumbnail) {
    return <Thumbnail file={file} scale={scale} />
  } else {
    return (
      <Box position="relative">
        <FontIcon file={file} scale={scale} />
        <HStack position="absolute" bottom="-5px" right="0px" spacing="2px">
          {file.isShared && <SharedBadge />}
          {file.ocr && <OcrBadge />}
        </HStack>
      </Box>
    )
  }
}

export default FileIcon

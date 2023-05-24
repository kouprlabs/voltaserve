import { Box } from '@chakra-ui/react'
import { File } from '@/api/file'
import FontIcon from './font-icon'
import SharedSign from './shared-sign'
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
        {file.isShared && <SharedSign bottom="-5px" right="0px" />}
      </Box>
    )
  }
}

export default FileIcon

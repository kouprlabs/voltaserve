import { Box } from '@chakra-ui/react'
import { FaFileImage } from 'react-icons/fa'
import { File } from '@/api/file'
import SharedSign from './shared-sign'
import Thumbnail from './thumbnail'

type ImageIconProps = {
  file: File
  scale: number
}

const ICON_FONT_SIZE = 72

const ImageIcon = ({ file, scale }: ImageIconProps) => {
  if (file.thumbnail) {
    return <Thumbnail file={file} scale={scale} />
  } else {
    return (
      <Box position="relative">
        <FaFileImage fontSize={`${ICON_FONT_SIZE * scale}px`} />
        {file.isShared && <SharedSign bottom="-5px" right="0px" />}
      </Box>
    )
  }
}

export default ImageIcon

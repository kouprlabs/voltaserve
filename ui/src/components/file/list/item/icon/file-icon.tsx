import { Box, HStack } from '@chakra-ui/react'
import { File, SnapshotStatus } from '@/client/api/file'
import ErrorBadge from './error-badge'
import FontIcon from './font-icon'
import NewBadge from './new-badge'
import ProcessingBadge from './processing-badge'
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
          {file.status === SnapshotStatus.New && <NewBadge />}
          {file.status === SnapshotStatus.Processing && <ProcessingBadge />}
          {file.status === SnapshotStatus.Error && <ErrorBadge />}
        </HStack>
      </Box>
    )
  }
}

export default FileIcon

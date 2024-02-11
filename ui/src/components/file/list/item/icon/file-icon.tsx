import { useMemo } from 'react'
import { Box, HStack } from '@chakra-ui/react'
import { SnapshotStatus } from '@/client/api/file'
import { CommonItemProps } from '@/types/file'
import ErrorBadge from './error-badge'
import FontIcon from './font-icon'
import NewBadge from './new-badge'
import ProcessingBadge from './processing-badge'
import SharedBadge from './shared-badge'
import Thumbnail from './thumbnail'

type FileIconProps = CommonItemProps

const FileIcon = ({ file, scale, viewType }: FileIconProps) => {
  const { bottom, right } = useMemo(() => {
    if (viewType === 'grid') {
      return { bottom: '-5px', right: '0px' }
    } else {
      return { bottom: '-7px', right: '0px' }
    }
  }, [viewType])
  if (file.thumbnail) {
    return <Thumbnail file={file} scale={scale} />
  } else {
    return (
      <Box position="relative">
        <FontIcon file={file} scale={scale} />
        <HStack position="absolute" bottom={bottom} right={right} spacing="2px">
          {file.isShared ? <SharedBadge /> : null}
          {file.status === SnapshotStatus.New ? <NewBadge /> : null}
          {file.status === SnapshotStatus.Processing ? (
            <ProcessingBadge />
          ) : null}
          {file.status === SnapshotStatus.Error ? <ErrorBadge /> : null}
        </HStack>
      </Box>
    )
  }
}

export default FileIcon

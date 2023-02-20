import { useMemo, useState } from 'react'
import { Box, Skeleton, Image } from '@chakra-ui/react'
import { FaFileImage } from 'react-icons/fa'
import { File } from '@/api/file'
import variables from '@/theme/variables'
import { ItemSize } from '..'
import FileListItemSharedSign from './shared-sign'

type FileListItemImageIconProps = {
  file: File
  size: ItemSize
}

const invalidSizeError = 'Invalid item size'

const FileListItemImageIcon = ({ file, size }: FileListItemImageIconProps) => {
  const width = useMemo<number>(() => {
    if (size === ItemSize.Normal) {
      return 130
    }
    if (size === ItemSize.Large) {
      return 230
    }
    throw invalidSizeError
  }, [size])
  const height = useMemo(() => {
    if (size === ItemSize.Normal) {
      return 90
    }
    if (size === ItemSize.Large) {
      return 190
    }
    throw invalidSizeError
  }, [size])
  const fontSize = useMemo(() => {
    if (size === ItemSize.Normal) {
      return 72
    }
    if (size === ItemSize.Large) {
      return 150
    }
    throw invalidSizeError
  }, [size])
  const [isLoading, setIsLoading] = useState(true)

  if (file.snapshots[0]?.thumbnail) {
    return (
      <Box position="relative" width={width} height={height}>
        <Image
          src={file.snapshots[0]?.thumbnail}
          width={isLoading ? 0 : width}
          height={isLoading ? 0 : height}
          style={{
            objectFit: 'cover',
            width: isLoading ? 0 : width,
            height: isLoading ? 0 : height,
            borderRadius: variables.borderRadiusSm,
            visibility: isLoading ? 'hidden' : 'visible',
          }}
          alt={file.name}
          onLoad={() => setIsLoading(false)}
        />
        {isLoading && (
          <Skeleton
            width={width}
            height={height}
            borderRadius={variables.borderRadiusSm}
          />
        )}
        {file.isShared && <FileListItemSharedSign bottom="-5px" right="-5px" />}
      </Box>
    )
  } else {
    return <FaFileImage fontSize={fontSize} />
  }
}

export default FileListItemImageIcon

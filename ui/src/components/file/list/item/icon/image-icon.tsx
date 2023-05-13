import { useMemo, useState } from 'react'
import { Box, Skeleton, Image } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { FaFileImage } from 'react-icons/fa'
import { File } from '@/api/file'
import { ItemSize } from '..'
import SharedSign from './shared-sign'

type ImageIconProps = {
  file: File
  size: ItemSize
}

const ImageIcon = ({ file, size }: ImageIconProps) => {
  const width = useMemo<number>(() => {
    if (size === ItemSize.Normal) {
      return 130
    } else if (size === ItemSize.Large) {
      return 230
    } else {
      throw new Error(`Invalid item size: ${size}`)
    }
  }, [size])
  const height = useMemo(() => {
    if (size === ItemSize.Normal) {
      return 90
    } else if (size === ItemSize.Large) {
      return 190
    } else {
      throw new Error(`Invalid item size: ${size}`)
    }
  }, [size])
  const fontSize = useMemo(() => {
    if (size === ItemSize.Normal) {
      return 72
    } else if (size === ItemSize.Large) {
      return 150
    } else {
      throw new Error(`Invalid item size: ${size}`)
    }
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
        {file.isShared && <SharedSign bottom="-5px" right="-5px" />}
      </Box>
    )
  } else {
    return <FaFileImage fontSize={fontSize} />
  }
}

export default ImageIcon

import { useMemo, useState } from 'react'
import {
  Box,
  Center,
  Image,
  Skeleton,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { IconPlay, variables } from '@koupr/ui'
import { File } from '@/client/api/file'
import { getSizeWithAspectRatio } from '@/helpers/aspect-ratio'
import * as fileExtension from '@/helpers/file-extension'
import SharedSign from './shared-sign'

const MAX_WIDTH = 130
const MAX_HEIGHT = 130

export function getThumbnailWidth(file: File, scale: number): string {
  if (file.thumbnail) {
    const { width } = getSizeWithAspectRatio(
      file.thumbnail.width,
      file.thumbnail.height,
      MAX_WIDTH,
      MAX_HEIGHT
    )
    return `${width * scale}px`
  } else {
    return `${MAX_WIDTH * scale}px`
  }
}

export function getThumbnailHeight(file: File, scale: number): string {
  if (file.thumbnail) {
    const { height } = getSizeWithAspectRatio(
      file.thumbnail.width,
      file.thumbnail.height,
      MAX_WIDTH,
      MAX_HEIGHT
    )
    return `${height * scale}px`
  } else {
    return `${MAX_HEIGHT * scale}px`
  }
}

type ThumbnailProps = {
  file: File
  scale: number
}

const Thumbnail = ({ file, scale }: ThumbnailProps) => {
  const width = useMemo(() => getThumbnailWidth(file, scale), [scale, file])
  const height = useMemo(() => getThumbnailHeight(file, scale), [scale, file])
  const [isLoading, setIsLoading] = useState(true)
  const borderColor = useColorModeValue('gray.300', 'gray.700')
  const [borderColorDecoded] = useToken('colors', [borderColor])
  const isVideo = useMemo(
    () =>
      file.original?.extension &&
      fileExtension.isVideo(file.original.extension),
    [file.original]
  )
  return (
    <Box position="relative" width={width} height={height}>
      <Image
        src={file.thumbnail?.base64}
        width={isLoading ? 0 : width}
        height={isLoading ? 0 : height}
        style={{
          objectFit: 'cover',
          width: isLoading ? 0 : width,
          height: isLoading ? 0 : height,
          border: '1px solid',
          borderColor: borderColorDecoded,
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
      {isVideo && (
        <Center
          position="absolute"
          top="0px"
          left="0px"
          width={width}
          height={height}
          opacity={0.5}
        >
          <IconPlay fontSize="40px" color="white" />
        </Center>
      )}
      {file.isShared && <SharedSign bottom="-5px" right="-5px" />}
    </Box>
  )
}

export default Thumbnail

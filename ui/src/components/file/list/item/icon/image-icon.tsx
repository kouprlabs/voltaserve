import { useMemo, useState } from 'react'
import {
  Box,
  Skeleton,
  Image,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { FaFileImage } from 'react-icons/fa'
import { File } from '@/api/file'
import SharedSign from './shared-sign'
import { getThumbnailHeight, getThumbnailWidth } from './thumbnail-size'

type ImageIconProps = {
  file: File
  scale: number
}

const ICON_FONT_SIZE = 72

const ImageIcon = ({ file, scale }: ImageIconProps) => {
  const width = useMemo(() => getThumbnailWidth(file, scale), [scale, file])
  const height = useMemo(() => getThumbnailHeight(file, scale), [scale, file])
  const iconFontSize = useMemo(() => {
    return `${ICON_FONT_SIZE * scale}px`
  }, [scale])
  const [isLoading, setIsLoading] = useState(true)
  const borderColor = useColorModeValue('gray.300', 'gray.700')
  const [borderColorDecoded] = useToken('colors', [borderColor])

  if (file.thumbnail) {
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
        {file.isShared && <SharedSign bottom="-5px" right="-5px" />}
      </Box>
    )
  } else {
    return <FaFileImage fontSize={iconFontSize} />
  }
}

export default ImageIcon

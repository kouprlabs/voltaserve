import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'

type ImageInfoProps = {
  file: File
}

const ImageInfo = ({ file }: ImageInfoProps) => {
  if (!file.original?.image) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Image dimensions</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        {file.original.image.width}x{file.original.image.height}
      </StatNumber>
    </Stat>
  )
}

export default ImageInfo

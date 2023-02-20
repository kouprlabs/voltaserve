import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { File } from '@/api/file'
import variables from '@/theme/variables'

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

import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'

export type FileInfoImageProps = {
  file: File
}

const FileInfoImage = ({ file }: FileInfoImageProps) => {
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

export default FileInfoImage

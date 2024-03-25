import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'

type FileInfoExtensionProps = {
  file: File
}

const FileInfoExtension = ({ file }: FileInfoExtensionProps) => {
  if (!file.original) {
    return null
  }
  return (
    <Stat>
      <StatLabel>File type</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        <Badge>{file.original.extension}</Badge>
      </StatNumber>
    </Stat>
  )
}

export default FileInfoExtension

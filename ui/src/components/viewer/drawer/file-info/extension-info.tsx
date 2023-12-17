import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'

type ExtensionInfoProps = {
  file: File
}

const ExtensionInfo = ({ file }: ExtensionInfoProps) => {
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

export default ExtensionInfo

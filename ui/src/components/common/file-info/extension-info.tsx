import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { File } from '@/api/file'
import variables from '@/theme/variables'

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

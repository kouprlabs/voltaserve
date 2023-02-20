import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { File } from '@/api/file'
import variables from '@/theme/variables'
import prettyBytes from '@/helpers/pretty-bytes'

type SizeInfoProps = {
  file: File
}

const SizeInfo = ({ file }: SizeInfoProps) => {
  if (!file.original) {
    return null
  }
  return (
    <Stat>
      <StatLabel>File size</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        {prettyBytes(file.original.size)}
      </StatNumber>
    </Stat>
  )
}

export default SizeInfo

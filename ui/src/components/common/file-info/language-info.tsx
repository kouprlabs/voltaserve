import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File } from '@/client/api/file'

type LanguageInfoProps = {
  file: File
}

const LanguageInfo = ({ file }: LanguageInfoProps) => {
  if (!file.language) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Language</StatLabel>
      <StatNumber fontSize={variables.bodyFontSize}>
        <Badge>{file.language}</Badge>
      </StatNumber>
    </Stat>
  )
}

export default LanguageInfo
